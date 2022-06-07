package stdlib

import (
	"fmt"
	. "rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testFuncData = []string{"ab", "100", "1000000", "5000000.123456", "22"}
var testFuncContext = KeyBuilderContextArray{
	Elements: testFuncData,
}

func mockContext(args ...interface{}) KeyBuilderContext {
	s := make([]string, len(args))
	for idx, arg := range args {
		s[idx] = fmt.Sprintf("%v", arg)
	}
	return &KeyBuilderContextArray{
		Elements: s,
	}
}

func testExpression(t *testing.T, context KeyBuilderContext, expression string, expected string) {
	kb, err := NewStdKeyBuilder().Compile(expression)
	assert.NoError(t, err)
	assert.NotNil(t, kb)
	if kb != nil {
		ret := kb.BuildKey(context)
		assert.Equal(t, expected, ret)
	}
}

func TestCoalesce(t *testing.T) {
	testExpression(t,
		mockContext("", "a", "b"),
		"{coalesce {0}} {coalesce a b c} {coalesce {0} {2}}",
		" a b")
}

func TestSimpleFunction(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{hi {2}} {hf {3}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1,000,000 5,000,000.1235", key)
}

func TestBucketing(t *testing.T) {
	testContext := mockContext("ab", "cd", "123")
	kb, _ := NewStdKeyBuilder().Compile("{bucket {2} 10} is bucketed")
	key := kb.BuildKey(testContext)
	assert.Equal(t, "120 is bucketed", key)
	assert.Equal(t, 2, kb.StageCount())
}

func TestBucket(t *testing.T) {
	testExpression(t,
		mockContext("1000", "1200", "1234"),
		"{bucket {0} 1000} {bucket {1} 1000} {bucket {2} 1000} {bucket {2} 100}",
		"1000 1000 1000 1200")
	testExpression(t, mockContext(), "{bucket abc 100} {bucket 1}", "<BUCKET-ERROR> <ARGN>")
}

func TestExpBucket(t *testing.T) {
	testExpression(t, mockContext("123", "1234", "12345"),
		"{expbucket {0}} {expbucket {1}} {expbucket {2}}", "100 1000 10000")
}

func TestClamp(t *testing.T) {
	testExpression(t, mockContext("100", "200", "1000", "-10"),
		"{clamp {0} 50 200}-{clamp {1} 50 200}-{clamp {2} 50 200}-{clamp {3} 50 200}",
		"100-200-max-min")
}

func TestByteSize(t *testing.T) {
	testExpression(t, &testFuncContext, "{bytesize {2}}", "977 KB")
	testExpression(t, &testFuncContext, "{bytesize {2} 2}", "976.56 KB")
}

func TestIfStatement(t *testing.T) {
	testExpression(t, mockContext("abc", "q"),
		`{if {0} {1} efg} {if {0} abc} {if {not {0}} a b} {if "" a} {if "" a b}`,
		"q abc b  b")
	testExpression(t, mockContext("abc efg"), `{if {eq {0} "abc efg"} beq}`, "beq")
}

func TestComparisonEquality(t *testing.T) {
	testExpression(t, mockContext("123", "1234"),
		"{eq {0} 123} {eq {0} 1234} {not {eq {0} abc}} {neq 1 2} {neq 1 1}",
		"123  1 1 ")
}

func TestComparisonExpression(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{and {lt {2} 10000000} {gt {1} 50}} {or 1 {eq abc 123}} {or {eq abc 123} {eq qef agg}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1 1 ", key)
}

func TestNotExpression(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{not {and {lt {2} 10000000} {gt {1} 50}}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "", key)
}

func TestComparisonLtGtEqual(t *testing.T) {
	testExpression(t, mockContext(), "{gte 1 1} {gte 1 2} {gte 2 1} {lte 1 1} {lte 1 2} {lte 2 1}",
		"1  1 1 1 ")
}

func TestStringPrefixSufix(t *testing.T) {
	testExpression(t, mockContext(), "{prefix abc a} {suffix abc c} {prefix abc b} {suffix abc b}",
		"abc abc  ")
}

func TestTabulator(t *testing.T) {
	testExpression(t, mockContext(), "{tab a b} {tab a b c}", "a\tb a\tb\tc")
}

func TestArray(t *testing.T) {
	testExpression(t, mockContext("q"), "{$ {0} {1} 22}", "q\x00\x0022")
	testExpression(t, mockContext("q"), `{$ "{0} hi" 22}`, "q hi\x0022")
	testExpression(t, mockContext("q"), "{$ {0}}", "q")
}

func TestHumanize(t *testing.T) {
	testExpression(t, mockContext(), "{hi 12345} {hf 12345.123512} {hi abc} {hf abc}",
		"12,345 12,345.1235 <BAD-TYPE> <BAD-TYPE>")
}

func TestFormat(t *testing.T) {
	testExpression(t, mockContext(), `{format "%10s" abc}`, "       abc")
}

func TestLike(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "ab", key)
}

func TestArithmatic(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumi 1} {sumi 1 a} {sumi a 1} {sumi 1 1 b}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}

func TestArithmaticf(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf {1} {4}} {multf {1} 2} {divf {1} 2} {subf {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticfError(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{sumf 1} {sumf 1 a} {sumf a 1} {sumf 1 2 a}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}
