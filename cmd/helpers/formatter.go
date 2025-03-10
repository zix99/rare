package helpers

import (
	"rare/pkg/logger"
	"rare/pkg/multiterm/termformat"

	"github.com/urfave/cli/v2"
)

var FormatFlag = &cli.StringFlag{
	Name:    "format",
	Usage:   "Defines a format expression for displayed values",
	Aliases: []string{"fmt"},
}

func BuildFormatter(expr string) (termformat.Formatter, error) {
	if expr == "" {
		return termformat.Default, nil
	}

	return termformat.FromExpression(expr)
}

func BuildFormatterOrFail(expr string) termformat.Formatter {
	f, err := BuildFormatter(expr)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}
	return f
}
