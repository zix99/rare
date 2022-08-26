package sorting

import (
	"sort"
	"strconv"
)

type NameSorter Sorter[string]

// Sorters

func ByName(a, b string) bool {
	return a < b
}

func ByNameSmart(a, b string) bool {
	v0, err0 := strconv.ParseFloat(a, 64)
	v1, err1 := strconv.ParseFloat(b, 64)
	if err0 == nil && err1 == nil && v0 != v1 {
		return v0 < v1
	}
	return a < b
}

// Sort string methods

func SortStrings(arr []string, sorter NameSorter) {
	ws := wrappedSorter[string]{arr, sorter}
	sort.Sort(&ws)
}

func SortStringsBy[T any](arr []T, sorter NameSorter, extractor func(obj T) string) {
	ws := wrappedSorter[T]{arr, func(a, b T) bool {
		return sorter(extractor(a), extractor(b))
	}}
	sort.Sort(&ws)
}
