package testutil

import (
	"testing"
)

// true if tests are run with -race (via build flag in race.go)
var IsRaceMode = false

func AssertZeroAlloc(t *testing.T, bench func(b *testing.B)) {
	if testing.Short() {
		t.Skip("skip short tests")
	}
	if IsRaceMode {
		t.Skip("skip race mode")
	}

	br := testing.Benchmark(bench)

	if allocs := br.AllocsPerOp(); allocs > 0 {
		t.Errorf("had %d allocs", allocs)
	}
	if allocb := br.AllocedBytesPerOp(); allocb > 0 {
		t.Errorf("had %d alloc bytes", allocb)
	}
}

// Asserts two arrays/slices share the same memory block
func AssertSameMemory[T any](t *testing.T, arr0, arr1 []T) {
	if !IsSameMemory(arr0, arr1) {
		t.Error("Slices don't share underlying memory")
	}
}

// Asserts two arrays/slices do NOT share the same memory block
func AssertNotSameMemory[T any](t *testing.T, arr0, arr1 []T) {
	if IsSameMemory(arr0, arr1) {
		t.Error("Slices don't share underlying memory")
	}
}

// Check if 2 slices point to the same underlying memory block
// Looks at last element in capacity to test this to work with slices correctly
func IsSameMemory[T any](arr0, arr1 []T) bool {
	cap0, cap1 := cap(arr0), cap(arr1)
	if cap0 == 0 || cap1 == 0 {
		return false
	}
	return &arr0[0:cap0][cap0-1] == &arr1[0:cap1][cap1-1]
}
