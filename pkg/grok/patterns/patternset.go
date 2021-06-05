package patterns

import (
	"bufio"
	"io"
	"strings"
)

type PatternSet struct {
	patterns map[string]string
}

func NewPatternSet() *PatternSet {
	return &PatternSet{
		patterns: make(map[string]string),
	}
}

func (s *PatternSet) Lookup(p string) (string, bool) {
	ret, ok := s.patterns[strings.ToUpper(p)]
	return ret, ok
}

func (s *PatternSet) LoadPatternFile(r io.Reader) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			break
		}

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") { // comment
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		s.AddPattern(parts[0], parts[1])
	}
}

func (s *PatternSet) AddPattern(name, expr string) {
	s.patterns[name] = expr
}

func (s *PatternSet) Patterns() []string {
	ret := make([]string, 0, len(s.patterns))
	for name := range s.patterns {
		ret = append(ret, name)
	}
	return ret
}
