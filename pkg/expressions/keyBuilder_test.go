package expressions

import (
	"bytes"
	"testing"
	"text/template"

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
	assert.Equal(t, 1, len(kb.stages))
}

func TestSimpleReplacement(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is {1}")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is cd", key)
	assert.Equal(t, 3, len(kb.stages))
}

func TestUnterminatedReplacement(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is {123")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is ", key)
	assert.Equal(t, 2, len(kb.stages))
}

func TestEscapedString(t *testing.T) {
	kb := NewKeyBuilder().Compile("{0} is \\{1\\} cool")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "ab is {1} cool", key)
	assert.Equal(t, 2, len(kb.stages))
}

func TestBucketing(t *testing.T) {
	kb := NewKeyBuilder().Compile("{bucket {2} 10} is bucketed")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "120 is bucketed", key)
	assert.Equal(t, 2, len(kb.stages))
}

func TestDeepKeys(t *testing.T) {
	kb := NewKeyBuilder().Compile("{{1} b} is bucketed")
	key := kb.BuildKey(&testContext)
	assert.Equal(t, "<Err:{1}> is bucketed", key)
}

func BenchmarkSimpleReplacement(b *testing.B) {
	kb := NewKeyBuilder().Compile("{0} is awesome")
	for n := 0; n < b.N; n++ {
		kb.BuildKey(&testContext)
	}
}

func BenchmarkGoTextTemplate(b *testing.B) {
	kb, _ := template.New("test").Parse("{a} is awesome")
	for n := 0; n < b.N; n++ {
		var buf bytes.Buffer
		kb.Execute(&buf, nil)
	}
}
