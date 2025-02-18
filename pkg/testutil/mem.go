package testutil

import (
	"testing"
)

func ZeroAlloc(t *testing.T, bench func(b *testing.B)) {
	br := testing.Benchmark(bench)
	if allocs := br.AllocsPerOp(); allocs > 0 {
		t.Errorf("had %d allocs", allocs)
	}
	if allocb := br.AllocedBytesPerOp(); allocb > 0 {
		t.Errorf("had %d alloc bytes", allocb)
	}
}

func SameMemory[T any](t *testing.T, arr0, arr1 []T) {
	if !IsSameMemory(arr0, arr1) {
		t.Error("Slices don't share underlying memory")
	}
}

func NotSameMemory[T any](t *testing.T, arr0, arr1 []T) {
	if IsSameMemory(arr0, arr1) {
		t.Error("Slices don't share underlying memory")
	}
}

func IsSameMemory[T any](arr0, arr1 []T) bool {
	cap0, cap1 := cap(arr0), cap(arr1)
	if cap0 == 0 || cap1 == 0 {
		return false
	}
	return &arr0[0:cap0][cap0-1] == &arr1[0:cap1][cap1-1]
}
