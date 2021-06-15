package stdlib

import (
	"testing"
)

func TestRepeatCharacter(t *testing.T) {
	testExpression(t,
		mockContext("4"),
		"{repeat a 2} {repeat b {0}}",
		"aa bbbb")
}

func TestAddingColor(t *testing.T) {
	testExpression(t,
		mockContext("what what"),
		"{color red {0}}",
		"what what")
}

func TestBarGraph(t *testing.T) {
	testExpression(t,
		mockContext(),
		"{bar 2 5 5}",
		"██")
}
