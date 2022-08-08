package multiterm

import (
	"io"
)

type VirtualTerm struct {
	lines []string
}

var _ MultilineTerm = &VirtualTerm{}

func NewVirtualTerm() *VirtualTerm {
	return &VirtualTerm{
		lines: make([]string, 0),
	}
}

func (s *VirtualTerm) WriteForLine(line int, text string) {
	if line >= len(s.lines) {
		old := s.lines
		s.lines = make([]string, line+1)
		copy(s.lines, old)
	}
	s.lines[line] = text
}

func (s *VirtualTerm) Close() {}

func (s *VirtualTerm) Get(line int) string {
	if line >= len(s.lines) || line < 0 {
		return ""
	}
	return s.lines[line]
}

func (s *VirtualTerm) LineCount() int {
	return len(s.lines)
}

// WriteToOutput writes to a terminal, preventing any potential wrapping
func (s *VirtualTerm) WriteToOutput(out io.Writer) {
	newLineBytes := []byte{'\n'}
	for _, line := range s.lines {
		WriteLineNoWrap(out, line)
		out.Write(newLineBytes)
	}
}
