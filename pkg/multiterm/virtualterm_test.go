package multiterm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualTerm(t *testing.T) {
	vt := NewVirtualTerm()
	vt.WriteForLine(0, "Hello")
	vt.WriteForLine(2, "Thar")

	assert.Equal(t, "Hello", vt.Get(0))
	assert.Equal(t, "", vt.Get(1))
	assert.Equal(t, 3, vt.LineCount())

	// Out of bounds
	assert.Equal(t, "", vt.Get(-1))
	assert.Equal(t, "", vt.Get(3))

	// Full write
	var sb strings.Builder
	vt.WriteToOutput(&sb)

	assert.Equal(t, "Hello\n\nThar\n", sb.String())
}
