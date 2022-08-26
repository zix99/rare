package sorting

import "sort"

type Sorter[T any] func(a, b T) bool

// Wrapped sorter (so sort package doesn't have to reflect in sort.Slice)
// Slight performance increase

type wrappedSorter[T any] struct {
	arr  []T
	less func(a, b T) bool
}

var _ sort.Interface = &wrappedSorter[string]{}

func (s *wrappedSorter[T]) Len() int {
	return len(s.arr)
}

func (s *wrappedSorter[T]) Swap(i, j int) {
	s.arr[i], s.arr[j] = s.arr[j], s.arr[i]
}

func (s *wrappedSorter[T]) Less(i, j int) bool {
	return s.less(s.arr[i], s.arr[j])
}

// Sorting helpers

func SortBy[T any](arr []T, sorter Sorter[T]) {
	ws := wrappedSorter[T]{arr, sorter}
	sort.Sort(&ws)
}

func Reverse[T ~func(a, b Q) bool, Q any](sorter T) T {
	return func(a, b Q) bool {
		return !sorter(a, b)
	}
}
