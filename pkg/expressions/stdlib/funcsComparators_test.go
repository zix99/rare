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

func TestUnlessStatement(t *testing.T) {
	testExpression(t, mockContext("abc"), `{unless {1} {0}} {unless abc efg} {unless "" bob}`, "abc  bob")
	testExpressionErr(t, mockContext("abc"), `{unless joe}`, "<ARGN>", ErrArgCount)
}

func TestComparisonEquality(t *testing.T) {
	testExpression(t, mockContext("123", "1234"),
		"{eq {0} 123} {eq {0} 1234} {not {eq {0} abc}} {neq 1 2} {neq 1 1}",
		"1  1 1 ")
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
}

func TestLike(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(mockContext("ab"))
	assert.Equal(t, "ab", key)
}
