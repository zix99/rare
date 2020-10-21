package multiterm

import (
	"fmt"
	"io"
)

type VirtualTerm struct {
	lines   map[int]string
	maxLine int
}

var _ MultilineTerm = &VirtualTerm{}

func NewVirtualTerm() *VirtualTerm {
	return &VirtualTerm{
		lines:   make(map[int]string),
		maxLine: 0,
	}
}

func (s *VirtualTerm) WriteForLine(line int, text string) {
	s.lines[line] = text
	if line > s.maxLine {
		s.maxLine = line
	}
}

func (s *VirtualTerm) Close() {}

func (s *VirtualTerm) Get(line int) string {
	return s.lines[line]
}

func (s *VirtualTerm) LineCount() int {
	return s.maxLine + 1
}

func (s *VirtualTerm) WriteToOutput(out io.Writer) {
	for i := 0; i <= s.maxLine; i++ {
		if l, ok := s.lines[i]; ok {
			WriteLineNoWrap(out, l)
		}
		fmt.Print("\n")
	}
}
