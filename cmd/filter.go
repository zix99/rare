package cmd

import (
	"fmt"

	. "rare/cmd/helpers"
	"rare/pkg/color"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	writeLines := c.Bool("line")
	customExtractor := c.IsSet("extract")

	extractor := BuildExtractorFromArguments(c)
	readChan := extractor.ReadChan()
	for {
		matchBatch, more := <-readChan
		if !more {
			break
		}
		for _, match := range matchBatch {
			if writeLines {
				fmt.Printf("%d: ", match.LineNumber)
			}
			if !customExtractor {
				fmt.Println(color.WrapIndices(match.Line, match.Indices[2:]))
			} else {
				fmt.Println(match.Extracted)
			}
		}
	}
	WriteExtractorSummary(extractor)
	return nil
}

// FilterCommand Exported command
func FilterCommand() *cli.Command {
	return AdaptCommandForExtractor(cli.Command{
		Name:      "filter",
		Usage:     "Filter incoming results with search criteria, and output raw matches",
		ShortName: "f",
		Action:    filterFunction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "line,l",
				Usage: "Output line numbers",
			},
		},
	})
}
