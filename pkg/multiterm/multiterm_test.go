package multiterm

import (
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

func (s *VirtualTerm) WriteForLine(line int, text string) {
	s.lines[line] = text
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
