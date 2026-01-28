package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zix99/rare/pkg/aggregation/sorting"
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/stringSplitter"

	"github.com/urfave/cli/v2"
)

var DefaultSortFlag = &cli.StringFlag{
	Name:  "sort",
	Usage: "Sorting method for display in format `key:order`. Keys: (v)alue, (t)ext, (s)mart, (n)umeric, (c)ontextual, (d)ate; Orders: (a)scending, (d)escending, (r)everse",
	Value: "smart",
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
		return nil, fmt.Errorf("unknown sort %s: %w", name, err)
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
	reverse = (realname == "value" || realname == "val" || realname == "v") // Value defaults descending

	if modifier, hasModifier := splitter.NextOk(); hasModifier {
		switch strings.ToLower(modifier) {
		case "rev", "reverse", "r":
			reverse = !reverse
		case "desc", "descending", "d":
			reverse = true
		case "asc", "ascending", "a":
			reverse = false
		default:
			return "", false, errors.New("invalid sort modifier, options: reverse, descending, ascending")
		}
	}

	return
}

func lookupSorter(name string) (sorting.NameValueSorter, error) {
	name = strings.ToLower(name)
	switch name {
	case "text", "t", "":
		return sorting.ValueNilSorter(sorting.ByName), nil
	case "numeric", "n":
		return sorting.ValueNilSorter(sorting.ByNumberStrict), nil
	case "smart", "s":
		return sorting.ValueNilSorter(sorting.ByNameSmart), nil
	case "contextual", "context", "c":
		return sorting.ValueNilSorter(sorting.ByContextual()), nil
	case "date", "d":
		return sorting.ValueNilSorter(sorting.ByDateWithContextual()), nil
	case "value", "val", "v":
		return sorting.ValueSorterEx(sorting.ByName), nil
	}
	return nil, errors.New("unknown sort, options: text, numeric, smart, contextual, date, value")
}
