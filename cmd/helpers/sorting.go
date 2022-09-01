package helpers

import (
	"errors"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/logger"
	"strings"

	"github.com/urfave/cli/v2"
)

func SortFlagWithDefault(dflt string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "sort",
		Usage: "Sets sorting method (value, text, numeric, contextual)",
		Value: dflt,
	}
}

// Defaults by-value
var SortFlag = SortFlagWithDefault("value")

var SortReverseFlag = &cli.BoolFlag{
	Name:    "sort-reverse",
	Aliases: []string{"reverse"},
	Usage:   "Reverses the display sort-order",
}

func AddSortFlag(command *cli.Command, defaultMode string) {
	if _, err := lookupSorter(defaultMode); err != nil {
		panic(err)
	}

	command.Flags = append(command.Flags,
		&cli.StringFlag{
			Name:  "sort",
			Usage: "Sets sorting method (value, text, numeric, contextual)",
			Value: defaultMode,
		},
		&cli.BoolFlag{
			Name:    "sort-reverse",
			Aliases: []string{"reverse"},
			Usage:   "Reverses the display sort-order",
		},
	)
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
	case "value":
		return sorting.ValueSorterEx(sorting.ByName), nil
	}
	return nil, errors.New("unknown sort")
}

func BuildSorter(name string, reverse bool) sorting.NameValueSorter {
	sorter, err := lookupSorter(name)
	if err != nil {
		logger.Fatalf("Unknown sort: %s (%v)", name, err)
		return nil
	}
	if reverse {
		sorter = sorting.Reverse(sorter)
	}
	return sorter
}

func BuildSorterFromFlags(c *cli.Context) sorting.NameValueSorter {
	return BuildSorter(c.String("sort"), c.Bool("sort-reverse"))
}
