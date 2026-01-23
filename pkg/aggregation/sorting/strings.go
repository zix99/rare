package sorting

import (
	"strconv"
)

type NameSorter Sorter[string]

// Sorters

// String compare
func ByName(a, b string) bool {
	return a < b
}

// Check if numeric, otherwise string compare
func ByNameSmart(a, b string) bool {
	v0, err0 := strconv.ParseFloat(a, 64)
	v1, err1 := strconv.ParseFloat(b, 64)
	if err0 == nil && err1 == nil {
		return v0 < v1
	}
	return a < b
}

// Use the first occurring number in a string as a comparer, otherwise string compare
func ByFirstNumber(a, b string) bool {
	v0, err0 := extractFirstNumber(a)
	v1, err1 := extractFirstNumber(b)
	if err0 == nil && err1 == nil {
		return v0 < v1
	}
	return a < b
}

func extractFirstNumber(s string) (float64, error) {
	num := ""
	found := false
	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			num += string(r)
			found = true
		} else if found {
			break
		}
	}
	if num == "" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseFloat(num, 64)
}
