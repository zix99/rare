package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalStore(t *testing.T) {
	glob := 3

	SwitchGlobal(&glob, 1)
	assert.Equal(t, 1, glob)
	SwitchGlobal(&glob, 2)
	assert.Equal(t, 2, glob)
	RevertGlobals()

	assert.Equal(t, 3, glob)
}
