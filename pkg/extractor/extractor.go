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
	Extracted   string
	LineNumber  uint64
	MatchNumber uint64
}

type Config struct {
	Posix   bool
	Regex   string
	Extract string
	Workers int
}

type Extractor struct {
	ReadChan     chan *Match
	regex        *regexp.Regexp
	readLines    uint64
	matchedLines uint64
	config       Config
	keyBuilder   *expressions.CompiledKeyBuilder
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

// async safe
func (s *Extractor) processLineSync(line string) {
	lineNum := atomic.AddUint64(&s.readLines, 1)
	matches := s.regex.FindAllStringSubmatch(line, -1)

	// Extract and forward to the ReadChan if there are matches
	if len(matches) > 0 {
		matchNum := atomic.AddUint64(&s.matchedLines, 1)
		context := expressions.KeyBuilderContextArray{
			Elements: matches[0],
		}
		s.ReadChan <- &Match{
			Line:        line,
			Groups:      matches[0],
			Extracted:   s.keyBuilder.BuildKey(&context),
			LineNumber:  lineNum,
			MatchNumber: matchNum,
		}
	}
}

// Create an extractor from an input channel
func New(input chan string, config *Config) *Extractor {
	extractor := Extractor{
		ReadChan:   make(chan *Match, 5),
		regex:      buildRegex(config.Regex, config.Posix),
		keyBuilder: expressions.NewKeyBuilder().Compile(config.Extract),
		config:     *config,
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
		close(extractor.ReadChan)
	}()

	return &extractor
}

func (s *Config) getWorkerCount() int {
	if s.Workers <= 0 {
		return 2
	}
	return s.Workers
}
