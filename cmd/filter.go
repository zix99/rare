package cmd

import (
	"fmt"

	"rare/pkg/color"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	writeLines := c.Bool("line")
	customExtractor := c.IsSet("extract")

	extractor := buildExtractorFromArguments(c)
	readChan := extractor.ReadChan()
	for {
		match, more := <-readChan
		if !more {
			break
		}
		if writeLines {
			fmt.Printf("%d: ", match.LineNumber)
		}
		if !customExtractor {
			fmt.Println(color.WrapIndices(match.Line, match.Indices[2:]))
		} else {
			fmt.Println(match.Extracted)
		}
	}
	writeExtractorSummary(extractor)
	return nil
}

// FilterCommand Exported command
func FilterCommand() *cli.Command {
	return &cli.Command{
		Name:      "filter",
		Usage:     "Filter incoming results with search criteria, and output raw matches",
		Action:    filterFunction,
		ArgsUsage: "<-|filename|glob...>",
		Flags: buildExtractorFlags(
			cli.BoolFlag{
				Name:  "line,l",
				Usage: "Output line numbers",
			},
		),
	}
}
