package multiterm

import "testing"

func TestSimpleTable(t *testing.T) {
	table := NewTable(5, 5)
	table.WriteRow(0, "a", "b", "c", "d")
	table.WriteRow(4, "a", "b", "c", "d")
	table.WriteRow(10, "a", "b", "c", "d")
}
