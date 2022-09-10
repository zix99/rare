package sorting

import (
	"time"

	"github.com/araddon/dateparse"
)

func ByDate(fallbackSort NameSorter) NameSorter {
	format := ""
	fallback := false

	return func(a, b string) bool {
		if !fallback {
			if format == "" {
				var err error
				if format, err = dateparse.ParseFormat(a); err != nil {
					fallback = true
				}
			}

			if format != "" {
				d0, err0 := time.Parse(format, a)
				d1, err1 := time.Parse(format, b)
				if err0 == nil && err1 == nil {
					return d0.Before(d1)
				} else {
					fallback = true
				}
			}
		}

		return fallbackSort(a, b)
	}
}

func ByDateWithContextual() NameSorter {
	return ByDate(ByContextual())
}
