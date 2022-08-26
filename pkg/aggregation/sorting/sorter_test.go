package sorting

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedSorter(t *testing.T) {
	ws := wrappedSorter[int]{
		arr:  []int{5, 2, 3, 1},
		less: func(a, b int) bool { return a < b },
	}
	assert.Equal(t, 4, ws.Len())
	assert.True(t, ws.less(1, 2))
	ws.Swap(1, 2)
	assert.Equal(t, []int{5, 3, 2, 1}, ws.arr)

	sort.Sort(&ws)
	assert.Equal(t, []int{1, 2, 3, 5}, ws.arr)
}

func TestSortBy(t *testing.T) {
	arr := []int{5, 4, 1, 2, 3}
	SortBy(arr, func(a, b int) bool { return a < b })
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr)
}

func TestReverseSort(t *testing.T) {
	arr := []int{5, 4, 1, 2, 3}
	SortBy(arr, Reverse(func(a, b int) bool { return a < b }))
	assert.Equal(t, []int{5, 4, 3, 2, 1}, arr)
}

// BenchmarkExtractSort-4   	 3838752	       329.0 ns/op	      64 B/op	       2 allocs/op
func BenchmarkExtractSort(b *testing.B) {
	type wrappedStruct struct {
		s string
	}
	list := []wrappedStruct{
		{"b"},
		{"c"},
		{"d"},
		{"e"},
		{"f"},
	}
	for i := 0; i < b.N; i++ {
		SortStringsBy(list, ByName, func(obj wrappedStruct) string { return obj.s })
	}
}
