package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeArray(t *testing.T) {
	assert.Equal(t, "abc", MakeArray("abc"))
	assert.Equal(t, "abc\x00def", MakeArray("abc", "def"))
}
