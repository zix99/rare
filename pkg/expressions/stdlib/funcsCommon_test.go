package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesce(t *testing.T) {
	testExpression(t,
		mockContext("", "a", "b"),
		"{coalesce {0}} {coalesce a b c} {coalesce {0} {2}}",
		" a b")
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
	testExpression(t, mockContext(), "{bucket abc 100} {bucket 1}", "<BAD-TYPE> <ARGN>")
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
