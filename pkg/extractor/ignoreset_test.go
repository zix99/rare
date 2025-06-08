package extractor

import (
	"testing"

	"github.com/zix99/rare/pkg/expressions"

	"github.com/stretchr/testify/assert"
)

func mockArrayContext(elements ...string) expressions.KeyBuilderContext {
	return &expressions.KeyBuilderContextArray{
		Elements: elements,
	}
}

func TestEmptyIgnoreSet(t *testing.T) {
	is, err := NewIgnoreExpressions()
	assert.NoError(t, err)
	assert.Nil(t, is)
}

func TestSimpleIgnoreSet(t *testing.T) {
	is, err := NewIgnoreExpressions("{eq {0} ignoreme}")
	assert.NoError(t, err)
	assert.True(t, is.IgnoreMatch(mockArrayContext("ignoreme")))
	assert.False(t, is.IgnoreMatch(mockArrayContext("notme")))
}
