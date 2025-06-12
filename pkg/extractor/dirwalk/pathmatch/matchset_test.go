package pathmatch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmtpyMatchSet(t *testing.T) {
	ms := MatchSet{}
	assert.False(t, ms.Matches("bla"))
}

func TestValidMatchSet(t *testing.T) {
	ms, err := NewMatchSet("*.go", "*.tmp")
	assert.NoError(t, err)

	assert.True(t, ms.Matches("bla.go"))
	assert.True(t, ms.Matches("test.tmp"))
	assert.False(t, ms.Matches("file.txt"))
}

func TestInvalidMatchSet(t *testing.T) {
	ms, err := NewMatchSet("[unclosed")
	assert.Error(t, err)
	assert.Nil(t, ms)
}
