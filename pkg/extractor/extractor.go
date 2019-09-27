package extractor

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type Match struct {
	Line   string
	Groups []string
}

type Extractor struct {
	ReadChan chan *Match
	regex    *regexp.Regexp
}

type Config struct {
	Posix bool
	Regex string
}

func buildRegex(s string, posix bool) *regexp.Regexp {
	if posix {
		return regexp.MustCompilePOSIX(s)
	}
	return regexp.MustCompile(s)
}

func (s *Extractor) processLineSync(line string) {
	matches := s.regex.FindAllStringSubmatch(line, -1)

	if len(matches) > 0 {
		s.ReadChan <- &Match{
			Line:   line,
			Groups: matches[0],
		}
	}
}

func NewExtractor(input chan string, config *Config) *Extractor {
	extractor := Extractor{
		ReadChan: make(chan *Match, 5),
		regex:    buildRegex(config.Regex, config.Posix),
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
