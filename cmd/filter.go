package cmd

import (
	"fmt"

	"rare/cmd/helpers"
	"rare/pkg/color"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	var (
		writeLines      = c.Bool("line")
		customExtractor = c.IsSet("extract")
	)

	batcher := helpers.BuildBatcherFromArguments(c)
	extractor := helpers.BuildExtractorFromArgumentsEx(c, batcher, "\t")

	readChan := extractor.ReadChan()
	for {
		matchBatch, more := <-readChan
		if !more {
			break
		}
		for _, match := range matchBatch {
			if writeLines {
				fmt.Printf("%s %s: ", color.Wrap(color.BrightGreen, match.Source), color.Wrapi(color.BrightYellow, match.LineNumber))
			}
			if !customExtractor {
				if len(match.Indices) == 2 {
					// Single match, highlight entire phrase
					fmt.Println(color.WrapIndices(match.Line, match.Indices))
				} else {
					// Multi-match groups, highlight individual groups
					fmt.Println(color.WrapIndices(match.Line, match.Indices[2:]))
				}
			} else {
				fmt.Println(match.Extracted)
			}
		}
	}
	helpers.WriteExtractorSummary(extractor)

	return helpers.DetermineErrorState(batcher, extractor, nil)
}

// FilterCommand Exported command
func filterCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:  "filter",
		Usage: "Filter incoming results with search criteria, and output raw matches",
		Description: `Filters incoming results by a regex, and output the match or an extracted expression.
		Unable to output contextual information due to the application's parallelism.  Use grep if you
		need that`,
		ShortName: "f",
		Action:    filterFunction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "line,l",
				Usage: "Output source file and line number",
			},
		},
	})
}
