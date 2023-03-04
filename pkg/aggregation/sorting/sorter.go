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

// Sort an array that can be sorted by Sorter
func Sort[TElem any, TSort ~func(a, b TElem) bool](arr []TElem, sorter TSort) {
	ws := wrappedSorter[TElem]{arr, sorter}
	sort.Sort(&ws)
}

// Sort an array of elements, by a sub-element, than be sorted by T
func SortBy[TElem any, TBy any, TSort ~func(a, b TBy) bool](arr []TElem, sorter TSort, extractor func(obj TElem) TBy) {
	ws := wrappedSorter[TElem]{arr, func(a, b TElem) bool {
		return sorter(extractor(a), extractor(b))
	}}
	sort.Sort(&ws)
}

// Reverse a sorter (`not` the comparer)
func Reverse[TElem any, TSort ~func(a, b TElem) bool](sorter TSort) TSort {
	return func(a, b TElem) bool {
		return !sorter(a, b)
	}
}
