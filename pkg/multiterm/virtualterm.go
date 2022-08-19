package multiterm

import (
	"io"
)

type VirtualTerm struct {
	lines  []string
	closed bool
}

var _ MultilineTerm = &VirtualTerm{}

func NewVirtualTerm() *VirtualTerm {
	return NewVirtualTermEx(0, 10)
}

func NewVirtualTermEx(size, cap int) *VirtualTerm {
	return &VirtualTerm{
		lines: make([]string, size, cap),
	}
}

func (s *VirtualTerm) WriteForLine(line int, text string) {
	if s.closed {
		panic("virtualterm closed")
	}

	for line >= len(s.lines) {
		s.lines = append(s.lines, "")
	}

	s.lines[line] = text
}

// Close the virtual term. Doesn't ever really need to close, but useful in testing
func (s *VirtualTerm) Close() {
	s.closed = true
}

// IsClosed checks if Close() was called
func (s *VirtualTerm) IsClosed() bool {
	return s.closed
}

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
