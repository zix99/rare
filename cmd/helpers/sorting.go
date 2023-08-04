package helpers

import (
	"errors"
	"fmt"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/logger"
	"rare/pkg/stringSplitter"
	"strings"

	"github.com/urfave/cli/v2"
)

var DefaultSortFlag = &cli.StringFlag{
	Name:  "sort",
	Usage: "Sorting method for display (value, text, numeric, contextual, date)",
	Value: "numeric",
}

// Create a sort flag with a different default value
func DefaultSortFlagWithDefault(dflt string) *cli.StringFlag {
	if _, err := lookupSorter(dflt); err != nil {
		panic(err)
	}

	flag := *DefaultSortFlag
	flag.Value = dflt
	return &flag
}

func BuildSorterOrFail(fullName string) sorting.NameValueSorter {
	sorter, err := BuildSorter(fullName)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}
	return sorter
}

func BuildSorter(fullName string) (sorting.NameValueSorter, error) {
	name, reverse, err := parseSort(fullName)
	if err != nil {
		return nil, fmt.Errorf("error parsing sort: %v", err)
	}

	sorter, err := lookupSorter(name)
	if err != nil {
		return nil, fmt.Errorf("unknown sort: %s", name)
	}
	if reverse {
		sorter = sorting.Reverse(sorter)
	}
	return sorter, nil
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
			return "", false, errors.New("invalid sort modifier")
		}
	}

	return
}

func lookupSorter(name string) (sorting.NameValueSorter, error) {
	name = strings.ToLower(name)
	switch name {
	case "text", "":
		return sorting.ValueNilSorter(sorting.ByName), nil
	case "numeric":
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
