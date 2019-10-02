package extractor

import (
	"bufio"
	"io"
	"regexp"
	"strings"
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
}

type Extractor struct {
	ReadChan     chan *Match
	regex        *regexp.Regexp
	readLines    uint64
	matchedLines uint64
	config       Config
}

func buildRegex(s string, posix bool) *regexp.Regexp {
	if posix {
		return regexp.MustCompilePOSIX(s)
	}
	return regexp.MustCompile(s)
}

func (s *Extractor) processLineSync(line string) {
	s.readLines++
	matches := s.regex.FindAllStringSubmatch(line, -1)

	// Extract and forward to the ReadChan if there are matches
	if len(matches) > 0 {
		s.matchedLines++
		s.ReadChan <- &Match{
			Line:        line,
			Groups:      matches[0],
			Extracted:   buildStringFromGroups(matches[0], s.config.Extract),
			LineNumber:  s.readLines,
			MatchNumber: s.matchedLines,
		}
	}
}

// Create an extractor from an input channel
func NewExtractor(input chan string, config *Config) *Extractor {
	extractor := Extractor{
		ReadChan: make(chan *Match, 5),
		regex:    buildRegex(config.Regex, config.Posix),
		config:   *config,
	}

	go func() {
		for {
			s, more := <-input
			if !more {
				break
			}
			extractor.processLineSync(s)
		}
		close(extractor.ReadChan)
	}()

	return &extractor
}

// Create an extractor for an io.Reader
func NewExtractorReader(reader io.Reader, config *Config) *Extractor {
	input := make(chan string)

	bufReader := bufio.NewReader(reader)
	go func() {
		for {
			line, err := bufReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			input <- strings.TrimSuffix(line, "\n")
		}
		close(input)
	}()

	return NewExtractor(input, config)
}
