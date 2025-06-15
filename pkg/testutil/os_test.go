package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWindows(t *testing.T) {
	SkipWindows(t)
	assert.False(t, IsWindows())
}
