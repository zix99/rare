package helpers

import (
	"errors"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/logger"
	"rare/pkg/stringSplitter"
	"strings"

	"github.com/urfave/cli/v2"
)

var DefaultSortFlag = &cli.StringFlag{
	Name:  "sort",
	Usage: "Sorting method for display (value, text, numeric, contextual, date)",
	Value: "smart",
}

func DefaultSortFlagWithDefault(dflt string) *cli.StringFlag {
	if _, err := lookupSorter(dflt); err != nil {
		panic(err)
	}

	flag := *DefaultSortFlag
	flag.Value = dflt
	return &flag
}

func BuildSorterOrFail(fullName string) sorting.NameValueSorter {
	name, reverse, err := parseSort(fullName)
	if err != nil {
		logger.Fatalf("Error parsing sort: %s", err)
		return nil
	}

	sorter, err := lookupSorter(name)
	if err != nil {
		logger.Fatalf("Unknown sort: %s", name)
		return nil
	}
	if reverse {
		sorter = sorting.Reverse(sorter)
	}
	return sorter
}

func parseSort(name string) (realname string, reverse bool, err error) {
	splitter := stringSplitter.Splitter{
		S:     name,
		Delim: ":",
	}

	realname = strings.ToLower(splitter.Next())
	reverse = (realname == "value") // Value defaults descending

	if modifier, hasModifier := splitter.NextOk(); hasModifier {
		switch strings.ToLower(modifier) {
		case "rev", "reverse":
			reverse = !reverse
		case "desc":
			reverse = true
		case "asc":
			reverse = false
		default:
			return "", false, errors.New("invalid modifier")
		}
	}

	return
}

func lookupSorter(name string) (sorting.NameValueSorter, error) {
	name = strings.ToLower(name)
	switch name {
	case "text", "":
		return sorting.ValueNilSorter(sorting.ByName), nil
	case "smart", "numeric":
		return sorting.ValueNilSorter(sorting.ByNameSmart), nil
	case "contextual", "context":
		return sorting.ValueNilSorter(sorting.ByContextual()), nil
	case "date":
		return sorting.ValueNilSorter(sorting.ByDateWithContextual()), nil
	case "value":
		return sorting.ValueSorterEx(sorting.ByName), nil
	}
	return nil, errors.New("unknown sort")
}
