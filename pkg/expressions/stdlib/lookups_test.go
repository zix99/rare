package stdlib

import "testing"

func TestLoadFile(t *testing.T) {
	testExpression(t, mockContext(), "{load ../../../cmd/testdata/graph.txt}", "bob 22\njack 93\njill 3\nmaria 19")
	testExpressionErr(t, mockContext(), "{load {0}}", "<CONST>", ErrConst)
}

func TestLookup(t *testing.T) {
	testExpression(t, mockContext("bob"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "33")
	testExpression(t, mockContext("bobert"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "")
	testExpression(t, mockContext("#test"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "22")
	testExpression(t, mockContext("#test"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt} \"#\"}", "")
}

func BenchmarkLoadFile(b *testing.B) {
	benchmarkExpression(b, mockContext(), "{load ../../../cmd/testdata/graph.txt}", "bob 22\njack 93\njill 3\nmaria 19")
}
