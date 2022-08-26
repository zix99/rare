package sorting

import (
	"testing"
)

// wrapped BenchmarkStringSort-4   	 6859735	       177.0 ns/op	      32 B/op	       1 allocs/op
func BenchmarkStringSort(b *testing.B) {
	list := []string{"b", "c", "d", "e", "f"}
	for i := 0; i < b.N; i++ {
		SortStrings(list, ByName)
	}
}
