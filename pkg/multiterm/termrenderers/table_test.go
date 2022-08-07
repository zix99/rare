package termrenderers

import (
	"rare/pkg/multiterm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	table := NewTable(multiterm.New(), 5, 5)
	table.WriteRow(0, "a", "b", "c", "d")
	table.WriteRow(4, "a", "b", "c", "d")
	table.WriteRow(10, "a", "b", "c", "d")
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
