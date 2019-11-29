package multiterm

import (
	"fmt"
	"testing"
)

// Term to use for testing
type VirtualTerm struct {
	lines map[int]string
}

func NewVirtualTerm() *VirtualTerm {
	return &VirtualTerm{
		lines: make(map[int]string),
	}
}

func (s *VirtualTerm) WriteForLine(line int, format string, args ...interface{}) {
	s.lines[line] = fmt.Sprintf(format, args...)
}

func (s *VirtualTerm) Close() {}

func (s *VirtualTerm) Get(line int) string {
	return s.lines[line]
}

func TestBasicMultiterm(t *testing.T) {
	mt := New()
	mt.WriteForLine(0, "Hello")
	mt.WriteForLine(1, "you")
	mt.WriteForLine(10, "There")
}
