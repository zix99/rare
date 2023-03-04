package expressions

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

func TestEmptyString(t *testing.T) {
	assert.Equal(t, []string{"a", "bc", ""}, splitTokenizedArguments(`a "bc" ""`))
	assert.Equal(t, []string{"a", "bc", "", "def"}, splitTokenizedArguments(`a bc "" def`))
}

func TestQuotedBraces(t *testing.T) {
	tokens := splitTokenizedArguments(`a "{bc\"" def`)
	assert.Equal(t, []string{"a", "{bc\"", "def"}, tokens)
}

func TestDeepQuoting(t *testing.T) {
	tokens := splitTokenizedArguments(`if {eq {0} "abc def"} abc`)
	assert.Equal(t, []string{"if", `{eq {0} "abc def"}`, "abc"}, tokens)

	tokens2 := splitTokenizedArguments(`eq {0} "abc def"`)
	assert.Equal(t, []string{"eq", "{0}", "abc def"}, tokens2)
}

func TestEscapeAtEnd(t *testing.T) {
	tokens := splitTokenizedArguments(`test a\`)
	assert.Equal(t, []string{"test", "a"}, tokens)
}

// BenchmarkArgSplitter-4   	 1119465	      1081 ns/op	     136 B/op	       6 allocs/op
func BenchmarkArgSplitter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		splitTokenizedArguments(`eq {0} "abc def"`)
	}
}
