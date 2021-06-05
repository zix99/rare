package grok

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPatternLookup struct{}

func (s *mockPatternLookup) Lookup(pattern string) (string, bool) {
	if pattern == "NOTEXIST" {
		return pattern, false
	}
	if pattern == "NESTED" {
		return "%{A:a} %{B:b}", true
	}
	return fmt.Sprintf("PAT:%s", pattern), true
}

var grok = NewEx(&mockPatternLookup{})

func TestNoPattern(t *testing.T) {
	testRewrite(t, "this is a test", "this is a test")
	testRewrite(t, `this is (\d+)`, "this is (\\d+)")
}

func TestEscapes(t *testing.T) {
	testRewrite(t, "Escaped %%", "Escaped %")
	testRewrite(t, "Escaped %%{expr}", "Escaped %{expr}")
	testRewrite(t, "Ending %", "Ending %")
	testRewrite(t, "No {expr}", "No {expr}")
}

func TestUnclosedExpression(t *testing.T) {
	_, err := grok.RewriteGrokPattern("this is a %{unclosed")
	assert.Error(t, err)
}

func TestSimplePattern(t *testing.T) {
	testRewrite(t, "Username: %{USERNAME}", "Username: PAT:USERNAME")
}

func TestNamedPattern(t *testing.T) {
	testRewrite(t, "Username: %{USERNAME:uname}", "Username: (?P<uname>PAT:USERNAME)")
}

func TestMissingPattern(t *testing.T) {
	_, err := grok.RewriteGrokPattern("%{NOTEXIST:abc}")
	assert.Error(t, err)
}

func TestTypedPattern(t *testing.T) {
	testRewrite(t, "%{INT:a:int}", "(?P<a>PAT:INT)")
}

func TestNestedPattern(t *testing.T) {
	testRewrite(t, "%{NESTED}", "(?P<a>PAT:A) (?P<b>PAT:B)")
}

func testRewrite(t *testing.T, expr, expected string) {
	ret, err := grok.RewriteGrokPattern(expr)
	assert.NoError(t, err)
	assert.Equal(t, expected, ret)
}
