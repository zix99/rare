package multiterm

import (
	"testing"
)

func TestBufferedTerm(t *testing.T) {
	vt := NewBufferedTerm()
	vt.WriteForLine(0, "hello")
	vt.Close()
}
