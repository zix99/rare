package multiterm

import (
	"fmt"
	"io"
	"sort"
)

type VirtualTerm struct {
	lines map[int]string
}

func NewVirtualTerm() *VirtualTerm {
	return &VirtualTerm{
		lines: make(map[int]string),
	}
}

func (s *VirtualTerm) WriteForLine(line int, text string) {
	s.lines[line] = text
}

func (s *VirtualTerm) Close() {}

func (s *VirtualTerm) Get(line int) string {
	return s.lines[line]
}

func (s *VirtualTerm) LineCount() int {
	return len(s.lines)
}

func (s *VirtualTerm) WriteToOutput(out io.Writer) {
	keys := make([]int, 0, len(s.lines))
	for k := range s.lines {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, idx := range keys {
		WriteLineNoWrap(out, s.lines[idx])
		eraseRemainingLine()
		fmt.Print("\n")
	}
}
