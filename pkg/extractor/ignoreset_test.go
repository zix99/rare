package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyIgnoreSet(t *testing.T) {
	is, err := NewIgnoreExpressions()
	assert.NoError(t, err)
	assert.Nil(t, is)
}

func TestSimpleIgnoreSet(t *testing.T) {
	is, err := NewIgnoreExpressions("{eq {0} ignoreme}")
	assert.NoError(t, err)
	assert.True(t, is.IgnoreMatch("ignoreme"))
	assert.False(t, is.IgnoreMatch("notme"))
}
