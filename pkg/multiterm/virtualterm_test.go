package multiterm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualTerm(t *testing.T) {
	vt := NewVirtualTerm()
	vt.WriteForLine(0, "Hello")
	vt.WriteForLinef(2, "Thar %s", "bob")

	assert.Equal(t, "Hello", vt.Get(0))
	assert.Equal(t, "", vt.Get(1))
	assert.Equal(t, 3, vt.LineCount())

	// Out of bounds
	assert.Equal(t, "", vt.Get(-1))
	assert.Equal(t, "", vt.Get(3))

	// Full write
	var sb strings.Builder
	vt.WriteToOutput(&sb)

	assert.Equal(t, "Hello\n\nThar bob\n", sb.String())

	// And close
	vt.Close()
	assert.True(t, vt.IsClosed())
	assert.Panics(t, func() {
		vt.WriteForLine(0, "will panic")
	})
}
