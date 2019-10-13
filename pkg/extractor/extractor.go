package extractor

import (
	"rare/pkg/expressions"
	"regexp"
	"sync"
	"sync/atomic"
)

type Match struct {
	Line        string
	Groups      []string
	Indices     []int
	Extracted   string
	LineNumber  uint64
	MatchNumber uint64
}

// Config for the extractor
type Config struct {
	Posix   bool      // Posix parse regex
	Regex   string    // Regex to find matches
	Extract string    // Extract these values from regex (expression)
	Workers int       // Workers to parse regex
	Ignore  IgnoreSet // Ignore these truthy expressions
}

type Extractor struct {
	readChan     chan *Match
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
	return s.readLines
}

func (s *Extractor) MatchedLines() uint64 {
	return s.matchedLines
}

func (s *Extractor) IgnoredLines() uint64 {
	return s.ignoredLines
}

func (s *Extractor) ReadChan() <-chan *Match {
	return s.readChan
}

func indexToSlices(s string, indexMatches []int) []string {
	strings := make([]string, len(indexMatches)/2)
	for i := 0; i < len(indexMatches)/2; i++ {
		strings[i] = s[indexMatches[i*2]:indexMatches[i*2+1]]
	}
	return strings
}

// async safe
func (s *Extractor) processLineSync(line string) {
	lineNum := atomic.AddUint64(&s.readLines, 1)
	matches := s.regex.FindStringSubmatchIndex(line)

	// Extract and forward to the ReadChan if there are matches
	if len(matches) > 0 {
		slices := indexToSlices(line, matches)
		if s.ignore == nil || !s.ignore.IgnoreMatch(slices...) {
			matchNum := atomic.AddUint64(&s.matchedLines, 1)

			context := expressions.KeyBuilderContextArray{
				Elements: slices,
			}
			s.readChan <- &Match{
				Line:        line,
				Groups:      slices,
				Indices:     matches,
				Extracted:   s.keyBuilder.BuildKey(&context),
				LineNumber:  lineNum,
				MatchNumber: matchNum,
			}
		} else {
			atomic.AddUint64(&s.ignoredLines, 1)
		}
	}
}

// New an extractor from an input channel
func New(input <-chan string, config *Config) *Extractor {
	extractor := Extractor{
		readChan:   make(chan *Match, 5),
		regex:      buildRegex(config.Regex, config.Posix),
		keyBuilder: expressions.NewKeyBuilder().Compile(config.Extract),
		config:     *config,
		ignore:     config.Ignore,
	}

	var wg sync.WaitGroup

	for i := 0; i < config.getWorkerCount(); i++ {
		wg.Add(1)
		go func() {
			for {
				s, more := <-input
				if !more {
					break
				}
				extractor.processLineSync(s)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(extractor.readChan)
	}()

	return &extractor
}

func (s *Config) getWorkerCount() int {
	if s.Workers <= 0 {
		return 2
	}
	return s.Workers
}
