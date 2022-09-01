package sorting

type NameValuePair struct {
	Name  string
	Value int64
}

type NameValueSorter Sorter[NameValuePair]

func ValueSorterEx(fallback NameSorter) NameValueSorter {
	return func(a, b NameValuePair) bool {
		if a.Value == b.Value {
			return fallback(a.Name, b.Name)
		}
		return a.Value > b.Value
	}
}

func ValueNilSorter(sorter NameSorter) NameValueSorter {
	return func(a, b NameValuePair) bool {
		return sorter(a.Name, b.Name)
	}
}

var (
	NVValueSorter = ValueSorterEx(ByName)
	NVNameSorter  = ValueNilSorter(ByName)
	NVSmartSorter = ValueNilSorter(ByNameSmart)
)
