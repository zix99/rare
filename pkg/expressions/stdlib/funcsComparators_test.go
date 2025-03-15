package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfStatement(t *testing.T) {
	testExpression(t, mockContext("abc", "q"),
		`{if {0} {1} efg} {if {0} abc} {if {not {0}} a b} {if "" a} {if "" a b}`,
		"q abc b  b")
	testExpression(t, mockContext("abc efg"), `{if {eq {0} "abc efg"} beq}`, "beq")
	testExpression(t, mockContext(), `{if {eq "" ""} true false}`, "true")
	testExpression(t, mockContext(), `{if {eq "" "abc"} true false}`, "false")
	testExpression(t, mockContext(), `{if {neq "" "abc"} true false}`, "true")
	testExpression(t, mockContext(), `{if {neq "" ""} true false}`, "false")
}

func TestSwitch(t *testing.T) {
	testExpression(t, mockContext("a"), "{switch {eq {0} a} isa {eq {0} b} isb 1 null}", "isa")
	testExpression(t, mockContext("b"), "{switch {eq {0} a} isa {eq {0} b} isb 1 null}", "isb")
	testExpression(t, mockContext("c"), "{switch {eq {0} a} isa {eq {0} b} isb 1 null}", "null")
	testExpression(t, mockContext("c"), "{switch {eq {0} a} isa {eq {0} b} isb null}", "null")
	testExpression(t, mockContext("a"), "{switch {eq {0} a} isa {eq {0} b} isb 1}", "isa")
	testExpressionErr(t, mockContext("a"), "{switch {eq {0} a}}", "<ARGN>", ErrArgCount)
}

func TestUnlessStatement(t *testing.T) {
	testExpression(t, mockContext("abc"), `{unless {1} {0}} {unless abc efg} {unless "" bob}`, "abc  bob")
	testExpressionErr(t, mockContext("abc"), `{unless joe}`, "<ARGN>", ErrArgCount)
}

func TestComparisonEquality(t *testing.T) {
	testExpression(t, mockContext("123", "1234"),
		"{eq {0} 123} {eq {0} 1234} {not {eq {0} abc}} {neq 1 2} {neq 1 1}",
		"1  1 1 ")
	testExpressionErr(t, mockContext(), "{eq a}", "<ARGN>", ErrArgCount)
}

func TestComparisonExpression(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{and {lt {2} 10000000} {gt {1} 50}} {or 1 {eq abc 123}} {or {eq abc 123} {eq qef agg}}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "1 1 ", key)
}

func TestNotExpression(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{not {and {lt {2} 10000000} {gt {1} 50}}}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "", key)
}

func TestComparisonLtGtEqual(t *testing.T) {
	testExpression(t, mockContext(), "{gte 1 1} {gte 1 2} {gte 2 1} {lte 1 1} {lte 1 2} {lte 2 1}",
		"1  1 1 1 ")
	testExpressionErr(t, mockContext(), "{gt a}", "<ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext(1), "{gt 1 a}", "<BAD-TYPE>", ErrNum)
	testExpressionErr(t, mockContext(1), "{gt a 1}", "<BAD-TYPE>", ErrNum)
	testExpression(t, mockContext("a"), "{gt {0} 1}", "<BAD-TYPE>")
	testExpression(t, mockContext("a"), "{gt 1 {0}}", "<BAD-TYPE>")
}

func TestLike(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(mockContext("ab"))
	assert.Equal(t, "ab", key)
}

// BenchmarkComparison/{gt_{0}_2}-4         	14308363	        81.00 ns/op	       0 B/op	       0 allocs/op
func BenchmarkComparison(b *testing.B) {
	benchmarkExpression(b, mockContext("5"), "{gt {0} 2}", "1")
}
