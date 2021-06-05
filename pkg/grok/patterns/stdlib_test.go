package patterns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardLibLoading(t *testing.T) {
	lib := Stdlib()
	assert.Greater(t, len(lib.patterns), 10)
}

func TestLookupCommonExpression(t *testing.T) {
	lib := Stdlib()
	v, ok := lib.Lookup("USERNAME")
	assert.True(t, ok)
	assert.Equal(t, "[a-zA-Z0-9._-]+", v)
}

func TestLookupMiss(t *testing.T) {
	lib := Stdlib()
	_, ok := lib.Lookup("NO_EXIST")
	assert.False(t, ok)
}

func TestLookupHelper(t *testing.T) {
	v, ok := Lookup("USERNAME")
	assert.True(t, ok)
	assert.Equal(t, "[a-zA-Z0-9._-]+", v)
}
