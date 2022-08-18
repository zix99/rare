package termrenderers

import (
	"rare/pkg/multiterm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	table := NewTable(vt, 5, 5)
	table.WriteRow(0, "a", "b", "c", "d")
	table.WriteRow(4, "a", "b", "c", "d")
	table.WriteRow(10, "a", "b", "c", "d")

	assert.Equal(t, "a b c d ", vt.Get(0))
	assert.Equal(t, "", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
	assert.Equal(t, "", vt.Get(3))
	assert.Equal(t, "a b c d ", vt.Get(4))
	assert.Equal(t, "", vt.Get(5))
	assert.Equal(t, "", vt.Get(10))
}

func TestSimpleTableVirtual(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	table := NewTable(vt, 2, 2)
	table.WriteRow(0, "a", "b", "c")
	table.WriteRow(1, "q")
	table.WriteRow(2, "abc")

	assert.Equal(t, "a b ", vt.Get(0))
	assert.Equal(t, "q ", vt.Get(1))
	assert.Equal(t, "", vt.Get(2))
}

func TestTableUnicodeAlignment(t *testing.T) {
	vt := multiterm.NewVirtualTerm()
	table := NewTable(vt, 5, 5)
	table.WriteRow(0, "abc", "b", "c", "d", "")
	table.WriteRow(1, "a", "bc", "c", "d", "")
	table.WriteRow(2, "a", "b✥", "c", "✥", "")

	assert.Equal(t, "abc b  c d  ", vt.Get(0))
	assert.Equal(t, "a   bc c d  ", vt.Get(1))
	assert.Equal(t, "a   b✥ c ✥  ", vt.Get(2))
}
