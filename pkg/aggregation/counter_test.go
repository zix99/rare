package aggregation

import (
	"rare/pkg/aggregation/sorting"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAddition(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	items := val.Items()

	assert.Equal(t, 1, len(items), "Expected length 1")
	assert.Equal(t, "test", items[0].Name)
	assert.Equal(t, int64(1), items[0].Item.count)
}

func TestInOrderItems(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	val.Sample("abc")
	val.Sample("abc")
	val.Sample("test")
	val.Sample("abc")
	val.Sample("qq")

	items := val.ItemsSortedBy(2, sorting.NVValueSorter)

	assert.Equal(t, 2, len(items), "Expected top 2")
	assert.Equal(t, "abc", items[0].Name)
	assert.Equal(t, "test", items[1].Name)
}

func TestInOrderItemsByKey(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	val.Sample("abc")
	val.Sample("abc")
	val.Sample("test")
	val.Sample("abc")
	val.Sample("qq")
	val.Sample("qq\x002")
	val.Sample("qq\x00bad")

	items := val.ItemsSortedBy(3, sorting.ValueNilSorter(sorting.ByName))

	assert.Equal(t, 3, len(items))
	assert.Equal(t, 3, val.GroupCount())
	assert.Equal(t, int64(8), val.Total())
	assert.Equal(t, uint64(1), val.ParseErrors())
	assert.Equal(t, "abc", items[0].Name)
	assert.Equal(t, int64(3), items[0].Item.Count())
	assert.Equal(t, "qq", items[1].Name)
	assert.Equal(t, int64(3), items[1].Item.Count())
	assert.Equal(t, "test", items[2].Name)

	reverseSort := val.ItemsSortedBy(3, sorting.Reverse(sorting.ValueNilSorter(sorting.ByName)))
	assert.Equal(t, 3, len(reverseSort))
	assert.Equal(t, "test", reverseSort[0].Name)
}
