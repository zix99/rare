package stdlib

import (
	"rare/pkg/testutil"
	"testing"
)

func TestLoadFile(t *testing.T) {
	testExpression(t, mockContext(), "{load ../../../cmd/testdata/graph.txt}", "bob 22\njack 93\njill 3\nmaria 19")
	testExpressionErr(t, mockContext(), "{load {0}}", "<CONST>", ErrConst)
	testExpressionErr(t, mockContext(), "{load a b}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext(), "{load notarealfile.txt}", "<FILE>", ErrFile)

	testutil.SwitchGlobal(&DisableLoad, true)
	defer testutil.RestoreGlobals()
	testExpressionErr(t, mockContext(), "{load ../../../cmd/testdata/graph.txt}", "<FILE>", ErrFile)
}

func TestLookup(t *testing.T) {
	testExpression(t, mockContext("bob"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "33")
	testExpression(t, mockContext("bobert"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "")
	testExpression(t, mockContext("#test"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt}}", "22")
	testExpression(t, mockContext("#test"), "{lookup {0} {load ../../../cmd/testdata/lookup.txt} \"#\"}", "")
	testExpressionErr(t, mockContext(), "{lookup a}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext(), "{lookup a b c d}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext("fn"), "{lookup fn {0}}", "<CONST>", ErrConst)
}

func TestHasKey(t *testing.T) {
	testExpression(t, mockContext("bob"), "{haskey {0} {load ../../../cmd/testdata/lookup.txt}}", "1")
	testExpression(t, mockContext("nop"), "{haskey {0} {load ../../../cmd/testdata/lookup.txt}}", "")
	testExpression(t, mockContext("key"), "{haskey {0} {load ../../../cmd/testdata/lookup.txt}}", "1")
	testExpression(t, mockContext("#test"), "{haskey {0} {load ../../../cmd/testdata/lookup.txt}}", "1")
	testExpression(t, mockContext("#test"), "{haskey {0} {load ../../../cmd/testdata/lookup.txt} #}", "")
	testExpressionErr(t, mockContext(), "{haskey a}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext(), "{haskey a b c d}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext("fn"), "{haskey fn {0}}", "<CONST>", ErrConst)
}

func BenchmarkLoadFile(b *testing.B) {
	benchmarkExpression(b, mockContext(), "{load ../../../cmd/testdata/graph.txt}", "bob 22\njack 93\njill 3\nmaria 19")
}
