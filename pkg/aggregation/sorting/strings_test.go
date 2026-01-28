package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameSort(t *testing.T) {
	assert.True(t, ByName("a", "b"))
}

func TestNameSmartSort(t *testing.T) {
	assert.True(t, ByNameSmart("a", "b"))
	assert.True(t, ByNameSmart("0.0", "1.0"))
	assert.True(t, ByNameSmart("1", "02"))
	assert.True(t, ByNameSmart("1", "b"))

	assert.True(t, ByNameSmart("a 123 a", "a 123 b"))
}

func TestStrictNumSort(t *testing.T) {
	assert.True(t, ByNumberStrict("1", "2"))
	assert.True(t, ByNumberStrict("01", "2"))
	assert.True(t, ByNumberStrict("1", "a"))

	assert.False(t, ByNumberStrict("a", "2"))

	assert.True(t, ByNumberStrict("a", "b"))
	assert.False(t, ByNumberStrict("b", "a"))
}

func TestStrictNumSortSet(t *testing.T) {
	arr := []string{"b", "a", "5", "41", "7", "q"}
	Sort(arr, ByNumberStrict)
	assert.Equal(t, []string{"5", "7", "41", "a", "b", "q"}, arr)
}

func TestSortStrings(t *testing.T) {
	arr := []string{"b", "c", "a", "q"}
	Sort(arr, ByNameSmart)
	assert.Equal(t, []string{"a", "b", "c", "q"}, arr)
}

func TestSortStringsBy(t *testing.T) {
	type wrapper struct {
		s string
	}
	arr := []wrapper{
		{"b"},
		{"a"},
		{"c"},
	}
	SortBy(arr, ByName, func(w wrapper) string { return w.s })

	assert.Equal(t, "a", arr[0].s)
	assert.Equal(t, "b", arr[1].s)
	assert.Equal(t, "c", arr[2].s)
}

func TestExtractNumber(t *testing.T) {
	v, ok := extractNumber("nonum")
	assert.Equal(t, int64(0), v)
	assert.False(t, ok)

	v, ok = extractNumber("")
	assert.Equal(t, int64(0), v)
	assert.False(t, ok)

	v, ok = extractNumber("123")
	assert.Equal(t, int64(123), v)
	assert.True(t, ok)

	v, ok = extractNumber("its 123")
	assert.Equal(t, int64(123), v)
	assert.True(t, ok)

	v, ok = extractNumber("its -123")
	assert.Equal(t, int64(-123), v)
	assert.True(t, ok)

	v, ok = extractNumber("its 0 or 123")
	assert.Equal(t, int64(0), v)
	assert.True(t, ok)

	v, ok = extractNumber("its 123-")
	assert.Equal(t, int64(123), v)
	assert.True(t, ok)

	v, ok = extractNumber("abc is 123 but another is 456")
	assert.Equal(t, int64(123), v)
	assert.True(t, ok)

	v, ok = extractNumber("abc -is 123 but another is 456")
	assert.Equal(t, int64(123), v)
	assert.True(t, ok)

	v, ok = extractNumber("abc is -123 but another is 456")
	assert.Equal(t, int64(-123), v)
	assert.True(t, ok)
}

// wrapped BenchmarkStringSort-4   	 6859735	       177.0 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStringSort(b *testing.B) {
	list := []string{"b", "c", "d", "e", "f"}
	for i := 0; i < b.N; i++ {
		Sort(list, ByName)
	}
}
