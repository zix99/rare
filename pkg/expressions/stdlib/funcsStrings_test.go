package stdlib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperLower(t *testing.T) {
	testExpression(t, mockContext("aBc"), "{upper {0}} {upper a b}", "ABC <ARGN>")
	testExpression(t, mockContext("aBc"), "{lower {0}} {lower a b}", "abc <ARGN>")
}

func TestSubstring(t *testing.T) {
	testExpression(t,
		mockContext("abcd"),
		"{substr {0} 0 2} {substr {0} 0 10} {substr {0} 3 2} {substr {0} 3 1} {substr 0}",
		"ab abcd d d <ARGN>")
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
		"{select {0} 0} {select {0} 1} {select {0} 2} {select {0} 3} {select 0} {select {1} 1}",
		"ab c d  <ARGN> q")
	testExpression(t, mockContext(), `{select "ab cd ef" 1}`, "cd")
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
