package sorting

import "strconv"

type NameSorter Sorter[string]

// Sorters

// String compare
func ByName(a, b string) bool {
	return a < b
}

// Check if numeric, or starts with number, otherwise string compare
func ByNameSmart(a, b string) bool {
	v0, ok0 := extractNumber(a)
	v1, ok1 := extractNumber(b)
	if ok0 && ok1 && v0 != v1 {
		return v0 < v1
	}
	return a < b
}

// Sorts numbers, always sorting numbers first. Any non-numbers will be sorted secondarily
func ByNumberStrict(a, b string) bool {
	v0, err0 := strconv.ParseFloat(a, 64)
	v1, err1 := strconv.ParseFloat(b, 64)
	if err0 == nil && err1 == nil {
		return v0 < v1
	}
	if err0 == nil && err1 != nil {
		return true
	}
	if err0 != nil && err1 == nil {
		return false
	}
	return a < b
}

// Extracts a number if it appears first in a string
func extractNumber(s string) (int64, bool) {
	var (
		total  = int64(0)
		negate = false
		found  = false
	)

	for _, r := range s {
		if r >= '0' && r <= '9' {
			total = total*10 + int64(r-'0')
			found = true
		} else if found {
			break
		} else if r == '-' {
			negate = true
		} else {
			negate = false
		}
	}

	if negate {
		total = -total
	}

	if found {
		return total, true
	}
	return 0, false
}
