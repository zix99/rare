package cmd

import (
	"fmt"
	"os"

	"rare/cmd/helpers"
	"rare/pkg/color"

	"github.com/urfave/cli/v2"
)

func filterFunction(c *cli.Context) error {
	var (
		writeLines      = c.Bool("line")
		customExtractor = c.IsSet("extract")
		numLineLimit    = uint64(c.Int64("num"))
		readLines       = uint64(0)
	)

	batcher := helpers.BuildBatcherFromArguments(c)
	extractor := helpers.BuildExtractorFromArgumentsEx(c, batcher, "\t")

	readChan := extractor.ReadFull()
OUTER_LOOP:
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

			readLines++
			if numLineLimit > 0 && readLines >= numLineLimit {
				break OUTER_LOOP
			}
		}
	}

	if numLineLimit > 0 {
		helpers.FWriteMatchSummary(os.Stderr, readLines, numLineLimit)
		os.Stderr.WriteString("\n")
	} else {
		helpers.WriteExtractorSummary(extractor)
	}

	return helpers.DetermineErrorState(batcher, extractor, nil)
}

// FilterCommand Exported command
func filterCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:  "filter",
		Usage: "Filter incoming results with search criteria, and output raw matches",
		Description: `Filters incoming results by a regex, and output the match of a single line
		or an extracted expression.`,
		Aliases:  []string{"f"},
		Action:   filterFunction,
		Category: cmdCatAnalyze,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "line",
				Aliases: []string{"l"},
				Usage:   "Output source file and line number",
			},
			&cli.Int64Flag{
				Name:    "num",
				Aliases: []string{"n"},
				Usage:   "Print the first NUM of lines seen (Not necessarily in-order)",
			},
		},
	})
}
