package sorting

import "strings"

type sortSet map[string]int

// But what about localization?!
// Unfortunately (or fortunately) golang only can encode dates in english, which helps for the time being

var weekdays = sortSet{
	"sunday":    0,
	"monday":    1,
	"tuesday":   2,
	"wednesday": 3,
	"thursday":  4,
	"friday":    5,
	"saturday":  6,

	"sun":   0,
	"mon":   1,
	"tue":   2,
	"tues":  2,
	"wed":   3,
	"thu":   4,
	"thur":  4,
	"thurs": 4,
	"fri":   5,
	"sat":   6,
}

var months = sortSet{
	"january":   0,
	"jan":       0,
	"february":  1,
	"feb":       1,
	"march":     2,
	"mar":       2,
	"april":     3,
	"apr":       3,
	"may":       4,
	"june":      5,
	"jun":       5,
	"july":      6,
	"jul":       6,
	"august":    7,
	"aug":       7,
	"september": 8,
	"sep":       8,
	"sept":      8,
	"october":   9,
	"oct":       9,
	"november":  10,
	"nov":       10,
	"december":  11,
	"dec":       11,
}

var sortSets = [...]sortSet{
	weekdays,
	months,
}

func ByContextualEx(fallbackSort NameSorter) NameSorter {
	var set sortSet
	fallback := false

	return func(a, b string) bool {
		if !fallback && set == nil {
			set = inferSortSetByValue(a)
			if set == nil {
				fallback = true
			}
		}

		// Try using the set
		if !fallback {
			lowerA := strings.ToLower(a)
			lowerB := strings.ToLower(b)
			v0, ok0 := set[lowerA]
			v1, ok1 := set[lowerB]
			if !ok0 || !ok1 {
				fallback = true
			} else {
				return v0 < v1
			}
		}

		// Fallback
		return fallbackSort(a, b)
	}
}

func ByContextual() NameSorter {
	return ByContextualEx(ByNameSmart)
}

func inferSortSetByValue(val string) sortSet {
	val = strings.ToLower(val)
	for _, set := range sortSets {
		if _, ok := set[val]; ok {
			return set
		}
	}
	return nil
}
