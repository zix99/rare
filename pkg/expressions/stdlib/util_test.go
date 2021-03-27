package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartialString(t *testing.T) {
	assert.True(t, isPartialString("m", "months"))
	assert.True(t, isPartialString("mo", "months"))
	assert.True(t, isPartialString("months", "months"))
	assert.False(t, isPartialString("t", "month"))
	assert.False(t, isPartialString("months", "month"))
}
