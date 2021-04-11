package aggregation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubKeyEmpty(t *testing.T) {
	sk := NewSubKeyCounter()
	assert.Equal(t, uint64(0), sk.ParseErrors())
	assert.Len(t, sk.SubKeys(), 0)
	assert.Len(t, sk.Items(), 0)
	assert.Len(t, sk.ItemsSorted(false), 0)
	assert.Len(t, sk.ItemsSorted(true), 0)
}

func TestSubKeyWithOnlyKeys(t *testing.T) {
	sk := NewSubKeyCounter()
	sk.SampleValue("test", "", 1)
	sk.SampleValue("test2", "", 3)
	sk.SampleValue("test2", "", 2)

	assert.Len(t, sk.SubKeys(), 1)
	assert.Len(t, sk.Items(), 2)

	items := sk.ItemsSorted(false)
	assert.Equal(t, "test", items[0].Name)
	assert.Equal(t, int64(1), items[0].Item.Count())
	assert.Equal(t, "test2", items[1].Name)
	assert.Equal(t, int64(5), items[1].Item.Count())
}

func TestSubKeyWithSubKeys(t *testing.T) {
	sk := NewSubKeyCounter()
	sk.SampleValue("test", "100", 1)
	sk.SampleValue("test", "200", 2)
	sk.SampleValue("test", "200", 2)
	sk.SampleValue("test2", "100", 3)

	assert.Len(t, sk.SubKeys(), 2)
	assert.Len(t, sk.Items(), 2)

	items := sk.ItemsSorted(false)
	assert.Len(t, items[0].Item.Items(), 2)
	assert.Len(t, items[1].Item.Items(), 2)
}

func TestSubKeyWithNullSample(t *testing.T) {
	sk := NewSubKeyCounter()
	sk.Sample("test")
	sk.Sample(fmt.Sprintf("%s\x00%s", "test", "abc"))
	sk.Sample(fmt.Sprintf("%s\x00%s\x00%d", "test", "abc", 5))

	assert.Len(t, sk.SubKeys(), 2)
	assert.Len(t, sk.Items(), 1)

	item := sk.Items()[0]
	assert.Equal(t, int64(7), item.Item.Count())
	assert.Equal(t, []int64{1, 6}, item.Item.Items())
}

func TestSubKeyParseError(t *testing.T) {
	sk := NewSubKeyCounter()
	sk.Sample(fmt.Sprintf("%s\x00%s\x00%s", "test", "test", "notnum"))

	assert.Len(t, sk.SubKeys(), 0)
	assert.Len(t, sk.Items(), 0)
	assert.Equal(t, uint64(1), sk.ParseErrors())
}
