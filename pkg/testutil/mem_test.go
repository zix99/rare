package testutil

import "testing"

func TestSameMemory(t *testing.T) {
	arr := make([]int, 5)
	arr0 := arr
	arr1 := arr[1:]
	arr2 := arr[1:2]
	arr3 := append(arr, 5, 6, 7)

	AssertSameMemory(t, arr, arr0)
	AssertSameMemory(t, arr, arr1)
	AssertSameMemory(t, arr, arr2)
	AssertNotSameMemory(t, arr, arr3)

	var blank []int
	AssertNotSameMemory(t, arr, blank)
}

func TestZeroAlloc(t *testing.T) {
	AssertZeroAlloc(t, BenchmarkEmpty)
}

func BenchmarkEmpty(b *testing.B) {}
