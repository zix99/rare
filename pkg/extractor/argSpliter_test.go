package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleToken(t *testing.T) {
	tokens := splitTokenizedArguments("abc")
	assert.Equal(t, []string{"abc"}, tokens)
}

func TestTokenList(t *testing.T) {
	tokens := splitTokenizedArguments("abc 1 qq")
	assert.Equal(t, []string{"abc", "1", "qq"}, tokens)
}

func TestQuotedToken(t *testing.T) {
	tokens := splitTokenizedArguments("abc 1 \"a b c\"")
	assert.Equal(t, []string{"abc", "1", "a b c"}, tokens)
}

func TestNestedToken(t *testing.T) {
	tokens := splitTokenizedArguments("abc {1} {a {2} 3}")
	assert.Equal(t, []string{"abc", "{1}", "{a {2} 3}"}, tokens)
}

func TestMultiSpaces(t *testing.T) {
	tokens := splitTokenizedArguments("a   b \"c \"  e")
	assert.Equal(t, []string{"a", "b", "c ", "e"}, tokens)
}

func TestTabSplit(t *testing.T) {
	tokens := splitTokenizedArguments("a\tb")
	assert.Equal(t, []string{"a", "b"}, tokens)
}

func TestEscaping(t *testing.T) {
	tokens := splitTokenizedArguments("a \\\"b c")
	assert.Equal(t, []string{"a", "\"b", "c"}, tokens)
}
