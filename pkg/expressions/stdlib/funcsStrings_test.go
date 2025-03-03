package stdlib

import (
	"rare/pkg/expressions"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	testExpression(t, mockContext("hello"), "{len {0}}", "5")
	testExpression(t, mockContext("hello"), "{len \"\"}", "0")
	testExpression(t, mockContext("hello"), "{len hi}", "2")
	testExpressionErr(t, mockContext("hello"), "{len {0} there}", "<ARGN>", ErrArgCount)
}

func TestUpperLower(t *testing.T) {
	testExpressionErr(t, mockContext("aBc"), "{upper {0}} {upper a b}", "ABC <ARGN>", ErrArgCount)
	testExpressionErr(t, mockContext("aBc"), "{lower {0}} {lower a b}", "abc <ARGN>", ErrArgCount)
}

func TestSubstring(t *testing.T) {
	testExpression(t,
		mockContext("abcd"),
		"{substr {0} 0 2} {substr {0} 0 10} {substr {0} 3 2} {substr {0} 3 1}",
		"ab abcd d d")
	testExpressionErr(t,
		mockContext("abcd"),
		"{substr 0}", "<ARGN>", ErrArgCount)
}

func TestSubstringOutOfBounds(t *testing.T) {
	testExpression(t,
		mockContext("abcd"),
		"{substr {0} -1 2} {substr {0} -2 2} {substr {0} -10 2} {substr {0} 3 4} {substr {0} 10 1}",
		"d cd ab d ")
}

func TestSelect(t *testing.T) {
	testExpression(t,
		mockContext("ab c d", "ab\tq"),
		"{select {0} 0} {select {0} 1} {select {0} 2} {select {0} 3} {select {1} 1}",
		"ab c d  q")
	testExpression(t, mockContext(), `{select "ab cd ef" 1}`, "cd")
	testExpressionErr(t, mockContext(), `{select 0}`, "<ARGN>", ErrArgCount)
}

func TestJoinEmpty(t *testing.T) {
	stage, err := kfJoin('-')([]expressions.KeyBuilderStage{})
	assert.NoError(t, err)
	assert.Equal(t, "", stage(mockContext()))
}

func TestSelectField(t *testing.T) {
	var s = "this  is\ta\ntest\x00really"
	assert.Equal(t, "this", selectField(s, 0))
	assert.Equal(t, "is", selectField(s, 1))
	assert.Equal(t, "a", selectField(s, 2))
	assert.Equal(t, "test", selectField(s, 3))
	assert.Equal(t, "really", selectField(s, 4))
	assert.Equal(t, "", selectField(s, 5))
}

func TestSelectFieldQuoted(t *testing.T) {
	assert.Equal(t, "a test", selectField(`this is "a test"`, 2))
	assert.Equal(t, "a test", selectField(`this is "a test" post`, 2))
	assert.Equal(t, "a test", selectField(`this " is" "a test"`, 2))
	assert.Equal(t, "  a test ", selectField(`this is "  a test "`, 2))
	assert.Equal(t, "  a test ", selectField(`this is "  a test `, 2))
}

func TestSimpleFunction(t *testing.T) {
	kb, _ := NewStdKeyBuilder().Compile("{hi {2}} {hf {3}}")
	key := kb.BuildKey(mockContext("ab", "100", "1000000", "5000000.123456", "22"))
	assert.Equal(t, "1,000,000 5,000,000.1235", key)
}

func TestPercentFunction(t *testing.T) {
	testExpression(t, mockContext("0.12345"), "{percent {0}}", "12.3%")
	testExpression(t, mockContext("0.12345"), "{percent {0} 2}", "12.35%")
	testExpressionErr(t, mockContext("0.12345"), "{percent {0} {0}}", "<CONST>", ErrConst)

	testExpression(t, mockContext("0.12345"), "{percent {0} 2 0.5}", "24.69%")
	testExpression(t, mockContext("50"), "{percent {0} 0 25 75}", "50%")

	testExpressionErr(t, mockContext(), "{percent 0 1 2 3 4 5}", "<ARGN>", ErrArgCount)
}

func TestDownscalers(t *testing.T) {
	testExpression(t, mockContext("1000000"), "{bytesize {0}}", "977 KB")
	testExpression(t, mockContext("1000000"), "{bytesize {0} 2}", "976.56 KB")
	testExpressionErr(t, mockContext("1000000"), "{bytesize {0} 2 3}", "<ARGN>", ErrArgCount)

	testExpression(t, mockContext("1000000"), "{bytesizesi {0}}", "1 mB")
	testExpression(t, mockContext("1000000"), "{bytesizesi {0} 2}", "1.00 mB")
	testExpressionErr(t, mockContext("1000000"), "{bytesizesi {0} 2 3}", "<ARGN>", ErrArgCount)

	testExpression(t, mockContext("5120000"), "{downscale {0}}", "5M")
	testExpression(t, mockContext("5120000"), "{downscale {0} 2}", "5.12M")
	testExpressionErr(t, mockContext("5120000"), "{downscale {0} 2 3}", "<ARGN>", ErrArgCount)
}

func TestFormat(t *testing.T) {
	testExpression(t, mockContext(), `{format "%10s" abc}`, "       abc")
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

func BenchmarkSplitFields(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strings.Fields("this  is\ta\ntest\x00really")
	}
}

func BenchmarkSelectItem(b *testing.B) {
	for i := 0; i < b.N; i++ {
		selectField("this  is\ta\ntest\x00really", 1)
	}
}
