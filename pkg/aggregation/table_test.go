package aggregation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleTable(t *testing.T) {
	table := NewTable(" ")
	table.Sample("b c")
	table.Sample("b c")
	table.Sample("a b")
	table.Sample("a c")
	table.Sample("b c")

	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())

	rows := table.OrderedRows()
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, "c", rows[0].Name())
	assert.Equal(t, "b", rows[1].Name())
	assert.Equal(t, int64(3), rows[0].Value("b"))
	assert.Equal(t, int64(1), rows[0].Value("a"))
}
