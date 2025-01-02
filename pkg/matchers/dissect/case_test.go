package dissect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexIgnoreCase(t *testing.T) {
	assert.Equal(t, 0, indexIgnoreCase("abc", "a"))
	assert.Equal(t, 0, indexIgnoreCase("abc", ""))

	assert.Equal(t, 1, indexIgnoreCase("abc", "bc"))
	assert.Equal(t, -1, indexIgnoreCase("abc", "ac"))

	assert.Equal(t, 0, indexIgnoreCase("abc", "abc"))

	assert.Equal(t, -1, indexIgnoreCase("abc", "bca"))
	assert.Equal(t, -1, indexIgnoreCase("abc", "abcd"))

	assert.Equal(t, 0, indexIgnoreCase("ABC", "a"))
	assert.Equal(t, 0, indexIgnoreCase("ABC", ""))

	assert.Equal(t, 1, indexIgnoreCase("ABC", "bc"))

	assert.Equal(t, 0, indexIgnoreCase("ABC", "abc"))

	assert.Equal(t, -1, indexIgnoreCase("ABC", "bca"))
	assert.Equal(t, -1, indexIgnoreCase("ABC", "abcd"))
}
