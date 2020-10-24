package stdlib

import "testing"

func TestSimpleCsv(t *testing.T) {
	testExpression(t,
		mockContext("a"),
		"{csv {0}}",
		"a")
	testExpression(t,
		mockContext("a", "b", "c"),
		"{csv {0} {1} {2}}",
		"a,b,c")
}

func TestEscapedCSV(t *testing.T) {
	testExpression(t,
		mockContext("a,", "b\"", "c"),
		"{csv {0} {1} {2}}",
		"\"a,\",\"b\"\"\",c")
}
