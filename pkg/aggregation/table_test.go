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
	table.Sample("b")

	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())

	rows := table.OrderedRows()
	assert.Equal(t, 2, len(rows))
	assert.Equal(t, "c", rows[0].Name())
	assert.Equal(t, "b", rows[1].Name())
	assert.Equal(t, int64(3), rows[0].Value("b"))
	assert.Equal(t, int64(1), rows[0].Value("a"))

	assert.Equal(t, 2, table.RowCount())
	assert.Equal(t, 2, table.ColumnCount())

	assert.Contains(t, table.Columns(), "a")
	assert.Contains(t, table.Columns(), "b")
	assert.Equal(t, []string{"b", "a"}, table.OrderedColumns())
	assert.Equal(t, uint64(1), table.ParseErrors())
}

func TestTableMultiIncrement(t *testing.T) {
	table := NewTable(" ")
	table.Sample("a b 1")
	table.Sample("b c 3")
	table.Sample("b c 3")
	table.Sample("b c -1")

	rows := table.OrderedRowsByName()
	assert.Equal(t, "c", rows[0].Name())
	assert.Equal(t, int64(5), rows[0].Value("b"))
	assert.Equal(t, "b", rows[1].Name())
	assert.Equal(t, int64(1), rows[1].Value("a"))
}
