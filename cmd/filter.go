package cmd

import (
	"bufio"
	"os"
	"unicode/utf8"

	"rare/cmd/helpers"
	"rare/pkg/color"
	"rare/pkg/logger"

	"github.com/urfave/cli/v2"
)

func filterFunction(c *cli.Context, fileGlobs ...string) error {
	var (
		writeLines      = c.Bool("line")
		customExtractor = c.IsSet("extract")
		onlyText        = c.Bool("text")
		numLineLimit    = uint64(c.Int64("num"))
		readLines       = uint64(0)
	)

	batcher := helpers.BuildBatcherFromArgumentsEx(c, fileGlobs...)
	extractor := helpers.BuildExtractorFromArgumentsEx(c, batcher, "\t")

	stdout := bufio.NewWriter(os.Stdout)

OUTER_LOOP:
	for matchBatch := range extractor.ReadFull() {
		for _, match := range matchBatch {
			if writeLines {
				color.WriteString(stdout, color.BrightGreen, match.Source)
				stdout.WriteByte(' ')
				color.WriteUint64(stdout, color.BrightYellow, match.LineNumber)
				stdout.WriteString(": ")
			}

			switch {
			case customExtractor:
				stdout.WriteString(match.Extracted)
			case onlyText && !utf8.ValidString(match.Line):
				color.WriteString(stdout, color.BrightBlue, "Binary Match")
			case len(match.Indices) == 2:
				// Single match, highlight entire phrase
				color.WrapIndices(stdout, match.Line, match.Indices)
			default:
				// Multi-match groups, highlight individual groups
				color.WrapIndices(stdout, match.Line, match.Indices[2:])
			}
			stdout.WriteByte('\n')

			readLines++
			if numLineLimit > 0 && readLines >= numLineLimit {
				break OUTER_LOOP
			}
		}

		// Flush after each batch to make file-following work as expected
		if err := stdout.Flush(); err != nil {
			logger.Fatal(helpers.ExitCodeInvalidUsage, err)
		}
	}

	// Final flush
	if err := stdout.Flush(); err != nil {
		logger.Fatal(helpers.ExitCodeInvalidUsage, err)
	}

	// Summary
	if numLineLimit > 0 {
		helpers.FWriteMatchSummary(os.Stderr, readLines, numLineLimit)
		os.Stderr.WriteString("\n")
	} else {
		helpers.WriteExtractorSummary(extractor)
	}

	return helpers.DetermineErrorState(batcher, extractor, nil)
}

func getFilterArgs(isSearch bool) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "line",
			Aliases: []string{"l"},
			Usage:   "Output source file and line number",
			Value:   isSearch,
		},
		&cli.Int64Flag{
			Name:    "num",
			Aliases: []string{"n"},
			Usage:   "Print the first NUM of lines seen (Not necessarily in-order)",
		},
		&cli.BoolFlag{
			Name:    "text",
			Aliases: []string{"a"},
			Usage:   "Only output lines that contain valid text",
			Value:   isSearch,
		},
	}
}

// FilterCommand Exported command
func filterCommand() *cli.Command {
	return helpers.AdaptCommandForExtractor(cli.Command{
		Name:  "filter",
		Usage: "Filter incoming results with search criteria, and output raw matches",
		Description: `Filters incoming results by a regex, and output the match of a single line
		or an extracted expression.`,
		Aliases: []string{"f"},
		Action: func(ctx *cli.Context) error {
			return filterFunction(ctx, ctx.Args().Slice()...)
		},
		Category: cmdCatAnalyze,
		Flags:    getFilterArgs(false),
	})
}

// Remap some arguments, and pass on to normal filter
func searchFunction(c *cli.Context) error {
	fileGlobs := c.Args().Slice()

	if len(fileGlobs) == 0 {
		logger.Fatal(helpers.ExitCodeInvalidUsage, "Missing match argument")
	}

	if !c.IsSet("match") && !c.IsSet("dissect") {
		c.Set("match", fileGlobs[0])
		fileGlobs = fileGlobs[1:]
	}

	if len(fileGlobs) == 0 {
		fileGlobs = append(fileGlobs, ".")
	}

	return filterFunction(c, fileGlobs...)
}

// Search command is very similar to filter, but with syntactic sugar to make
// it easier to discover things in a directory
func searchCommand() *cli.Command {
	command := helpers.AdaptCommandForExtractor(cli.Command{
		Name:        "search",
		Usage:       "Searches current directory recursively for a regex match",
		Description: `Same as filter, with defaults to easily search with a regex. Alias for: filter -IRla -m`,
		Action:      searchFunction,
		Category:    cmdCatAnalyze,
		Flags:       getFilterArgs(true),
	})

	command.ArgsUsage = "<regex> " + command.ArgsUsage

	// modify some defaults
	helpers.ModifyArgOrPanic(command, "recursive", func(flag *cli.BoolFlag) {
		flag.Value = true
	})
	helpers.ModifyArgOrPanic(command, "ignore-case", func(flag *cli.BoolFlag) {
		flag.Value = true
	})

	return command
}
