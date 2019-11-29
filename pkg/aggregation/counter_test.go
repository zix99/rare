package aggregation

import (
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

func collectChan(c chan interface{}) []interface{} {
	items := make([]interface{}, 0)
	for ele := range c {
		items = append(items, ele)
	}
	return items
}

func TestInOrderItems(t *testing.T) {
	val := NewCounter()
	val.Sample("test")
	val.Sample("abc")
	val.Sample("abc")
	val.Sample("test")
	val.Sample("abc")
	val.Sample("qq")

	items := val.ItemsTop(2)

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

	items := val.ItemsSortedByKey(3, false)

	assert.Equal(t, 3, len(items))
	assert.Equal(t, 3, val.GroupCount())
	assert.Equal(t, uint64(0), val.ParseErrors())
	assert.Equal(t, "abc", items[0].Name)
	assert.Equal(t, int64(3), items[0].Item.Count())
	assert.Equal(t, "qq", items[1].Name)
	assert.Equal(t, "test", items[2].Name)

	reverseSort := val.ItemsSortedByKey(3, true)
	assert.Equal(t, 3, len(reverseSort))
	assert.Equal(t, "test", reverseSort[0].Name)
}
