package extractor

import (
	"rare/pkg/expressions"
	"regexp"
	"sync"
	"sync/atomic"
	"unsafe"
)

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
	Indices    []int   // match indices as returned by regexp
	Extracted  string  // The extracted expression
	LineNumber uint64  // Line number
	Source     string  // Source name
}

// Config for the extractor
type Config struct {
	Posix   bool      // Posix parse regex
	Regex   string    // Regex to find matches
	Extract string    // Extract these values from regex (expression)
	Workers int       // Workers to parse regex
	Ignore  IgnoreSet // Ignore these truthy expressions
}

// Extractor is the representation of the reader
//  Expects someone to consume its ReadChan()
type Extractor struct {
	readChan     chan []Match
	regex        *regexp.Regexp
	readLines    uint64
	matchedLines uint64
	ignoredLines uint64
	config       Config
	keyBuilder   *expressions.CompiledKeyBuilder
	ignore       IgnoreSet
}

func buildRegex(s string, posix bool) *regexp.Regexp {
	if posix {
		return regexp.MustCompilePOSIX(s)
	}
	return regexp.MustCompile(s)
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

func (s *Extractor) ReadChan() <-chan []Match {
	return s.readChan
}

// async safe
func (s *Extractor) processLineSync(source string, lineNum uint64, line BString) (Match, bool) {
	atomic.AddUint64(&s.readLines, 1)
	matches := s.regex.FindSubmatchIndex(line)

	// Extract and forward to the ReadChan if there are matches
	if len(matches) > 0 {
		// Speed is more important here than safety
		// By default, casting to string will copy() data from bytes to
		//   a string instance, but we can safely point to the existing bytes
		//   as a pointer instead
		lineStringPtr := *(*string)(unsafe.Pointer(&line))
		expContext := SliceSpaceExpressionContext{
			linePtr: lineStringPtr,
			indices: matches,
		}
		if s.ignore == nil || !s.ignore.IgnoreMatch(&expContext) {
			extractedKey := s.keyBuilder.BuildKey(&expContext)

			if len(extractedKey) > 0 {
				atomic.AddUint64(&s.matchedLines, 1)
				return Match{
					bLine:      line,
					Line:       lineStringPtr,
					Indices:    matches,
					Extracted:  extractedKey,
					LineNumber: lineNum,
					Source:     source,
				}, true
			}

			atomic.AddUint64(&s.ignoredLines, 1)
		} else {
			atomic.AddUint64(&s.ignoredLines, 1)
		}
	}
	return Match{}, false
}

// New an extractor from an input channel
func New(inputBatch <-chan InputBatch, config *Config) (*Extractor, error) {
	compiledExpression, err := expressions.NewKeyBuilder().Compile(config.Extract)
	if err != nil {
		return nil, err
	}

	extractor := Extractor{
		readChan:   make(chan []Match, 5),
		regex:      buildRegex(config.Regex, config.Posix),
		keyBuilder: compiledExpression,
		config:     *config,
		ignore:     config.Ignore,
	}

	var wg sync.WaitGroup

	for i := 0; i < config.getWorkerCount(); i++ {
		wg.Add(1)
		go func() {
			for {
				batch, more := <-inputBatch
				if !more {
					break
				}

				var matchBatch []Match
				for idx, s := range batch.Batch {
					if match, ok := extractor.processLineSync(batch.Source, batch.BatchStart+uint64(idx), s); ok {
						if matchBatch == nil {
							// Initialize to expected cap (only if we have any matches)
							matchBatch = make([]Match, 0, len(batch.Batch))
						}
						matchBatch = append(matchBatch, match)
					}
				}
				if len(matchBatch) > 0 {
					extractor.readChan <- matchBatch
				}
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(extractor.readChan)
	}()

	return &extractor, nil
}

func (s *Config) getWorkerCount() int {
	if s.Workers <= 0 {
		return 2
	}
	return s.Workers
}
