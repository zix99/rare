package extractor

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/expressions/funclib"
	"github.com/zix99/rare/pkg/matchers"
	"github.com/zix99/rare/pkg/slicepool"
)

/*
Extractor is the main working pool that runs all CPU-intensive tasks (matcher and expressions)
Because of this, this area is extra-optimized, and so looks a little funny

Most aggregators will use ReadSimple(), which only returns the matches, and is significantly
faster than its ReadFull counterpart.  Because it's so sensitive, the code path is provided twice
rather than made generic
*/

// BString a []byte representation of a string (used for performance over string-copies)
type BString []byte

// InputBatch represents a batch of input
type InputBatch struct {
	Batch      []BString
	Source     string
	BatchStart uint64
}

// Match is a single given match
type Match struct {
	bLine      BString // Keep the pointer around next to line
	Line       string  // Unsafe pointer to bLine (no-copy)
	Indices    []int   // match indices as returned by matcher
	Extracted  string  // The extracted expression
	LineNumber uint64  // Line number
	Source     string  // Source name
}

// Config for the extractor
type Config struct {
	Matcher matchers.Factory // Matcher
	Extract string           // Extract these values from matcher (expression)
	Workers int              // Workers to parse matcher
	Ignore  IgnoreSet        // Ignore these truthy expressions
}

func (s *Config) getWorkerCount() int {
	if s.Workers <= 0 {
		return 2
	}
	return s.Workers
}

// Extractor is the representation of the reader
//
//	Expects someone to consume its ReadChan()
type Extractor struct {
	config     *Config
	keyBuilder *expressions.CompiledKeyBuilder
	ignore     IgnoreSet

	readLines    uint64
	matchedLines uint64
	ignoredLines uint64

	input <-chan InputBatch
}

func (s *Extractor) ReadLines() uint64 {
	return atomic.LoadUint64(&s.readLines)
}

func (s *Extractor) MatchedLines() uint64 {
	return atomic.LoadUint64(&s.matchedLines)
}

func (s *Extractor) IgnoredLines() uint64 {
	return atomic.LoadUint64(&s.ignoredLines)
}

func New(inputBatch <-chan InputBatch, config *Config) (*Extractor, error) {
	kb, err := funclib.NewKeyBuilder().Compile(config.Extract)
	if err != nil {
		return nil, err
	}

	ext := &Extractor{
		config:     config,
		keyBuilder: kb,
		ignore:     config.Ignore,
		input:      inputBatch,
	}

	return ext, nil
}

func (s *Extractor) workerFull(output chan<- []Match) {
	matcher := s.config.Matcher.CreateInstance()
	exprCtx := &SliceSpaceExpressionContext{
		nameTable: matcher.SubexpNameTable(),
	}
	matchBuf := make([]int, 0, matcher.MatchBufSize())
	bufCopyPool := slicepool.NewIntPool(1000 * matcher.MatchBufSize())

	for batch := range s.input {
		var (
			matchBatch   []Match
			ignoredCount int = 0
		)

		// setup
		atomic.AddUint64(&s.readLines, uint64(len(batch.Batch)))
		exprCtx.source = batch.Source

		// Process each line
		for idx, line := range batch.Batch {
			matches := matcher.FindSubmatchIndexDst(line, matchBuf)

			if len(matches) > 0 {
				// Speed is more important here than safety
				// By default, casting to string will copy() data from bytes to
				//   a string instance, but we can safely point to the existing bytes
				//   as a pointer instead
				lineStringPtr := *(*string)(unsafe.Pointer(&line))

				// A context is created for each "instance", and since a context isn't shared beyond building a key
				//   it's significantly faster to reuse a single context per goroutine
				exprCtx.linePtr = lineStringPtr
				exprCtx.indices = matches
				exprCtx.lineNum = batch.BatchStart + uint64(idx)

				// Check ignore, if possible
				if s.ignore == nil || !s.ignore.IgnoreMatch(exprCtx) {
					extractedKey := s.keyBuilder.BuildKey(exprCtx)

					// Extracted a key
					if len(extractedKey) > 0 {
						if matchBatch == nil {
							matchBatch = make([]Match, 0, len(batch.Batch))
						}

						// Need a copy of the local buffer (gets overwritten) to return
						matchCopy := bufCopyPool.Get(len(matches))
						copy(matchCopy, matches)

						matchBatch = append(matchBatch, Match{
							bLine:      line,
							Line:       lineStringPtr,
							Indices:    matchCopy,
							Extracted:  extractedKey,
							LineNumber: exprCtx.lineNum,
							Source:     batch.Source,
						})
					} else {
						ignoredCount++
					}

				} else {
					ignoredCount++
				}
			}
		}

		if ignoredCount > 0 {
			atomic.AddUint64(&s.ignoredLines, uint64(ignoredCount))
		}

		// Emit batch if there is data
		if len(matchBatch) > 0 {
			atomic.AddUint64(&s.matchedLines, uint64(len(matchBatch)))
			output <- matchBatch
		}
	}
}

// Read full match definitions, at the cost of lower performance (about 15% slower clock time / cpu time)
func (s *Extractor) ReadFull() <-chan []Match {
	return startWorkers(s.config.getWorkerCount(), s.workerFull)
}

func (s *Extractor) workerSimple(output chan<- []string) {
	matcher := s.config.Matcher.CreateInstance()
	exprCtx := &SliceSpaceExpressionContext{
		nameTable: matcher.SubexpNameTable(),
	}
	matchBuf := make([]int, 0, matcher.MatchBufSize())

	for batch := range s.input {
		var (
			matchBatch   []string
			ignoredCount int = 0
		)

		// setup
		atomic.AddUint64(&s.readLines, uint64(len(batch.Batch)))
		exprCtx.source = batch.Source

		// Process each line
		for idx, line := range batch.Batch {
			matches := matcher.FindSubmatchIndexDst(line, matchBuf)

			if len(matches) > 0 {
				// Speed is more important here than safety
				// By default, casting to string will copy() data from bytes to
				//   a string instance, but we can safely point to the existing bytes
				//   as a pointer instead
				lineStringPtr := *(*string)(unsafe.Pointer(&line))

				// A context is created for each "instance", and since a context isn't shared beyond building a key
				//   it's significantly faster to reuse a single context per goroutine
				exprCtx.linePtr = lineStringPtr
				exprCtx.indices = matches
				exprCtx.lineNum = batch.BatchStart + uint64(idx)

				// Check ignore, if possible
				if s.ignore == nil || !s.ignore.IgnoreMatch(exprCtx) {
					extractedKey := s.keyBuilder.BuildKey(exprCtx)

					// Extracted a key
					if len(extractedKey) > 0 {
						if matchBatch == nil {
							matchBatch = make([]string, 0, len(batch.Batch))
						}
						matchBatch = append(matchBatch, extractedKey)
					} else {
						ignoredCount++
					}

				} else {
					ignoredCount++
				}
			}
		}

		if ignoredCount > 0 {
			atomic.AddUint64(&s.ignoredLines, uint64(ignoredCount))
		}

		// Emit batch if there is data
		if len(matchBatch) > 0 {
			atomic.AddUint64(&s.matchedLines, uint64(len(matchBatch)))
			output <- matchBatch
		}
	}
}

func (s *Extractor) ReadSimple() <-chan []string {
	return startWorkers(s.config.getWorkerCount(), s.workerSimple)
}

func startWorkers[T string | Match](count int, worker func(output chan<- []T)) <-chan []T {
	var wg sync.WaitGroup
	output := make(chan []T, count*2)

	for range count {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(output)
		}()
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}
