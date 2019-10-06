package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []string{"ab", "cd", "123"}

type TestContext struct{}

func (s *TestContext) GetMatch(idx int) string {
	return testData[idx]
}

var testContext = TestContext{}

func TestSimpleKey(t *testing.T) {
	kb := NewKeyBuilder().Compile("test 123")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "test 123", key)
}

func TestSimpleReplacement(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is {1}")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is cd", key)
}

func TestUnterminatedReplacement(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is {123")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is ", key)
}

func TestEscapedString(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is \\{1\\} cool")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is {1} cool", key)
}

func TestBucketing(t *testing.T) {
	kb := NewKeyBuilder().Compile("{bucket 2 10} is bucketed")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "120 is bucketed", key)
}

func BenchmarkSimpleReplacement(b *testing.B) {
	kb := NewKeyBuilder().Compile("{0} is awesome")
	for n := 0; n < b.N; n++ {
		kb.BuildKey(&testContext)
	}
}
