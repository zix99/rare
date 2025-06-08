package stdlib

import (
	"testing"

	"github.com/zix99/rare/pkg/multiterm/termunicode"
	"github.com/zix99/rare/pkg/testutil"
)

func TestRepeatCharacter(t *testing.T) {
	testExpression(t,
		mockContext("4"),
		"{repeat a 2} {repeat b {0}}",
		"aa bbbb")
	testExpressionErr(t,
		mockContext("4"),
		"{repeat a} {repeat a a} {repeat {0} {0}}",
		"<ARGN> <BAD-TYPE> <CONST>",
		ErrArgCount, ErrConst)
}

func TestAddingColor(t *testing.T) {
	testExpression(t,
		mockContext("what what"),
		"{color red {0}}",
		"what what")
	testExpressionErr(t, mockContext("what what"), "{color bla {0}}", "<ENUM>", ErrEnum)
	testExpressionErr(t, mockContext("what what"), "{color {0} {0}}", "<CONST>", ErrConst)
	testExpressionErr(t,
		mockContext("what waht"),
		"{color a}", "<ARGN>", ErrArgCount)
}

func TestBarGraph(t *testing.T) {
	testutil.SwitchGlobal(&termunicode.UnicodeEnabled, false)
	defer testutil.RestoreGlobals()

	testExpression(t, mockContext(), "{bar 2 5 5}", "||")
	testExpression(t, mockContext(), "{bar 10 100 10}", "|")
	testExpression(t, mockContext(), "{bar 10 100 10 log10}", "|||||")
	testExpressionErr(t, mockContext(), "{bar 2 {0} 5}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext(), "{bar 10 100 10 badlog}", "<ENUM>", ErrEnum)
	testExpressionErr(t, mockContext(), "{bar 10 100 10 {0}}", "<CONST>", ErrConst)
}

func TestBarGraphErrorCases(t *testing.T) {
	testExpressionErr(t, mockContext(), "{bar 1}", "<ARGN>", ErrArgCount)
	testExpression(t, mockContext(), "{bar a 2 3}", "<BAD-TYPE>")
	testExpressionErr(t, mockContext(), "{bar 3 {1} {2}}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext(), "{bar 3 3 {2}}", "<BAD-TYPE>", ErrNum)
	testExpression(t, mockContext("a"), "{bar {0} 5 5}", "<BAD-TYPE>")
}
