package cmd

import (
	"fmt"

	"rare/pkg/color"
	. "rare/cmd/helpers"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	writeLines := c.Bool("line")
	customExtractor := c.IsSet("extract")

	extractor := BuildExtractorFromArguments(c)
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
	WriteExtractorSummary(extractor)
	return nil
}

// FilterCommand Exported command
func FilterCommand() *cli.Command {
	return &cli.Command{
		Name:      "filter",
		Usage:     "Filter incoming results with search criteria, and output raw matches",
		Action:    filterFunction,
		ArgsUsage: "<-|filename|glob...>",
		Flags: BuildExtractorFlags(
			cli.BoolFlag{
				Name:  "line,l",
				Usage: "Output line numbers",
			},
		),
	}
}
