package multiterm

import (
	"testing"
)

func TestBasicMultiterm(t *testing.T) {
	mt := New()
	mt.WriteForLine(0, "Hello")
	mt.WriteForLine(1, "you")
	mt.WriteForLine(10, "There")
}
