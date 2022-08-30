package sorting

import "sort"

type NameValue interface {
	SortName() string
	SortValue() int64
}

type NameValueSorter Sorter[NameValue]

func ValueSorter() NameValueSorter {
	return ValueSorterEx(ByName)
}

func ValueSorterEx(fallback NameSorter) NameValueSorter {
	return func(a, b NameValue) bool {
		v0, v1 := a.SortValue(), b.SortValue()
		if v0 == v1 {
			return fallback(a.SortName(), b.SortName())
		}
		return v0 > v1
	}
}

func ValueNilSorter(sorter NameSorter) NameValueSorter {
	return func(a, b NameValue) bool {
		return sorter(a.SortName(), b.SortName())
	}
}

var (
	ValueNameSorter  = ValueNilSorter(ByName)
	ValueSmartSorter = ValueNilSorter(ByNameSmart)
)

// SortNameValue interfaces.
// Needs to be separate because of the coersion into the interface itself
func SortNameValue[T NameValue](arr []T, sorter NameValueSorter) {
	ws := wrappedSorter[T]{arr, func(a, b T) bool {
		return sorter(a, b)
	}}
	sort.Sort(&ws)
}
