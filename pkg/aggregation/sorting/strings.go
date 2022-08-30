package sorting

import (
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
