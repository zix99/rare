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
	assert.True(t, ByNameSmart("1", "b"))
}

func TestSortStrings(t *testing.T) {
	arr := []string{"b", "c", "a", "q"}
	SortStrings(arr, ByNameSmart)
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
	SortStringsBy(arr, ByName, func(w wrapper) string { return w.s })

	assert.Equal(t, "a", arr[0].s)
	assert.Equal(t, "b", arr[1].s)
	assert.Equal(t, "c", arr[2].s)
}

// wrapped BenchmarkStringSort-4   	 6859735	       177.0 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStringSort(b *testing.B) {
	list := []string{"b", "c", "d", "e", "f"}
	for i := 0; i < b.N; i++ {
		SortStrings(list, ByName)
	}
}
