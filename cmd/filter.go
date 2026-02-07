package cmd

import (
	"bufio"
	"os"
	"os/signal"
	"unicode/utf8"

	"github.com/zix99/rare/cmd/helpers"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/logger"

	"github.com/urfave/cli/v2"
)

func filterFunction(c *cli.Context, fileGlobs ...string) error {
	var (
		writeLines      = c.Bool("line")
		customExtractor = c.IsSet("extract")
		onlyText        = c.Bool("text")
		summarize       = c.Bool("summary")
		numLineLimit    = uint64(c.Int64("num"))
		readLines       = uint64(0)
	)

	batcher, walker := helpers.BuildBatcherFromArgumentsEx(c, fileGlobs...)
	extractor := helpers.BuildExtractorFromArgumentsEx(c, batcher, "\t")

	stdout := bufio.NewWriter(os.Stdout)

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)
	interrupted := false

	readChan := extractor.ReadFull()

OUTER_LOOP:
	for {
		select {
		case <-exitSignal:
			interrupted = true
			break OUTER_LOOP
		case matchBatch, more := <-readChan:
			if !more {
				break OUTER_LOOP
			}

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
				logger.Fatal(helpers.ExitCodeOutputError, err)
			}
		}
	}

	// Final flush
	if err := stdout.Flush(); err != nil {
		logger.Fatal(helpers.ExitCodeOutputError, err)
	}

	// Summary
	if summarize {
		helpers.WriteBatcherSummary(os.Stderr, batcher, walker)
	}

	if numLineLimit > 0 {
		helpers.WriteMatchSummary(os.Stderr, readLines, numLineLimit)
	} else {
		os.Stderr.WriteString(helpers.BuildExtractorSummary(extractor, 0))
	}
	os.Stderr.WriteString("\n")

	return helpers.DetermineErrorState2(interrupted, batcher, extractor, nil)
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
		&cli.BoolFlag{
			Name:    "summary",
			Aliases: []string{"s"},
			Usage:   "Output a summary to stderr when done",
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
