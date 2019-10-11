package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	writeLines := c.Bool("line")

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		if writeLines {
			fmt.Printf("%d: %s\n", match.LineNumber, match.Extracted)
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
