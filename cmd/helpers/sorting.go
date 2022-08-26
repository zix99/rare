package helpers

import (
	"errors"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/logger"
	"strings"
)

func DecideSorterByName(name string) (sorting.NameValueSorter, error) {
	name = strings.ToLower(name)
	switch name {
	case "text", "":
		return sorting.ValueNilSorter(sorting.ByName), nil
	case "smart", "numeric":
		return sorting.ValueNilSorter(sorting.ByNameSmart), nil
	case "contextual", "context":
		return sorting.ValueNilSorter(sorting.ByContextual()), nil
	case "value":
		return sorting.ValueSorter(), nil
	}
	return nil, errors.New("unknown sort")
}

func BuildSorter(name string, reverse bool) sorting.NameValueSorter {
	sorter, err := DecideSorterByName(name)
	if err != nil {
		logger.Fatalf("Unknown sort: %s (%v)", name, err)
		return nil
	}
	if reverse {
		sorter = sorting.Reverse(sorter)
	}
	return sorter
}
