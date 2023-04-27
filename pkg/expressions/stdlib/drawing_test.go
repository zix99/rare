package stdlib

import (
	"testing"
)

func TestRepeatCharacter(t *testing.T) {
	testExpression(t,
		mockContext("4"),
		"{repeat a 2} {repeat b {0}}",
		"aa bbbb")
	testExpressionErr(t,
		mockContext("4"),
		"{repeat a} {repeat a a}",
		"<ARGN> <BAD-TYPE>",
		ErrArgCount)
}

func TestAddingColor(t *testing.T) {
	testExpression(t,
		mockContext("what what"),
		"{color red {0}}",
		"what what")
	testExpressionErr(t,
		mockContext("what waht"),
		"{color a}", "<ARGN>", ErrArgCount)
}

func TestBarGraph(t *testing.T) {
	testExpression(t,
		mockContext(),
		"{bar 2 5 5}",
		"██")
}

func TestBarGraphErrorCases(t *testing.T) {
	testExpressionErr(t, mockContext(), "{bar 1}", "<ARGN>", ErrArgCount)
	testExpression(t, mockContext(), "{bar a 2 3}", "<BAD-TYPE>")
	testExpressionErr(t, mockContext(), "{bar 3 {1} {2}}", "<BAD-TYPE>", ErrNum)
	testExpression(t, mockContext("a"), "{bar {0} 5 5}", "<BAD-TYPE>")
}
