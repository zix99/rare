package expressions

import (
	"fmt"
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
	kb, err := NewKeyBuilder().Compile(expression)
	assert.NoError(t, err)
	ret := kb.BuildKey(context)
	assert.Equal(t, expected, ret)
}

func TestCoalesce(t *testing.T) {
	testExpression(t,
		mockContext("", "a", "b"),
		"{coalesce a b c} {coalesce {0} {2}}",
		"a b")
}

func TestSimpleFunction(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{hi {2}} {hf {3}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1,000,000 5,000,000.1235", key)
}

func TestBucket(t *testing.T) {
	testExpression(t,
		mockContext("1000", "1200", "1234"),
		"{bucket {0} 1000} {bucket {1} 1000} {bucket {2} 1000} {bucket {2} 100}",
		"1000 1000 1000 1200")
	testExpression(t, mockContext(), "{bucket abc 100}", "<BUCKET-ERROR>")
}

func TestExpBucket(t *testing.T) {
	testExpression(t, mockContext("123", "1234", "12345"),
		"{expbucket {0}} {expbucket {1}} {expbucket {2}}", "100 1000 10000")
}

func TestByteSize(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{bytesize {2}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "976 KB", key)
}

func TestComparisonEquality(t *testing.T) {
	testExpression(t, mockContext("123", "1234"),
		"{eq {0} 123} {eq {0} 1234} {not {eq {0} abc}} {neq 1 2} {neq 1 1}",
		"123  1 1 ")
}

func TestComparisonExpression(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{and {lt {2} 10000000} {gt {1} 50}} {or 1 {eq abc 123}} {or {eq abc 123} {eq qef agg}}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "1 1 ", key)
}

func TestNotExpression(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{not {and {lt {2} 10000000} {gt {1} 50}}}")
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

func TestHumanize(t *testing.T) {
	testExpression(t, mockContext(), "{hi 12345} {hf 12345.123512} {hi abc} {hf abc}",
		"12,345 12,345.1235 <BAD-TYPE> <BAD-TYPE>")
}

func TestJson(t *testing.T) {
	testExpression(t, mockContext(`{"abc":123}`), `{json {0} abc}`, "123")
}

func TestFormat(t *testing.T) {
	testExpression(t, mockContext(), `{format "%10s" abc}`, "       abc")
}

func TestLike(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{like {0} \"a\"}{like {0} c}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "ab", key)
}

func TestArithmatic(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumi {1} {4}} {multi {1} 2} {divi {1} 2} {subi {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticError(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumi 1} {sumi 1 a} {sumi a 1} {sumi 1 1 b}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}

func TestArithmaticf(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumf {1} {4}} {multf {1} 2} {divf {1} 2} {subf {1} 10}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "122 200 50 90", key)
}

func TestArithmaticfError(t *testing.T) {
	kb, _ := NewKeyBuilder().Compile("{sumf 1} {sumf 1 a} {sumf a 1} {sumf 1 2 a}")
	key := kb.BuildKey(&testFuncContext)
	assert.Equal(t, "<ARGN> <BAD-TYPE> <BAD-TYPE> <BAD-TYPE>", key)
}
