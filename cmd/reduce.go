package cmd

import (
	"fmt"
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
	"rare/pkg/color"
	"rare/pkg/csv"
	"rare/pkg/expressions/stdlib"
	"rare/pkg/logger"
	"rare/pkg/multiterm/termrenderers"
	"strings"

	"github.com/urfave/cli/v2"
)

func reduceFunction(c *cli.Context) error {
	var (
		accumExprs     = c.StringSlice("accumulator")
		groupExpr      = c.StringSlice("group")
		defaultInitial = c.String("initial")
		table          = c.Bool("table")
		sort           = c.String("sort")
		sortReverse    = c.Bool("sort-reverse")
		rowCount       = c.Int("rows")
		colCount       = c.Int("cols")
	)

	vt := helpers.BuildVTermFromArguments(c)
	batcher := helpers.BuildBatcherFromArguments(c)
	extractor := helpers.BuildExtractorFromArguments(c, batcher)

	aggr := aggregation.NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	// Set up groups
	for _, group := range groupExpr {
		name, val := parseKeyValue(group)
		if err := aggr.AddGroupExpr(name, val); err != nil {
			logger.Fatalf(helpers.ExitCodeInvalidUsage, "Error compiling group expression %s: %s", group, err)
		}
	}

	// Set up expressions
	maxKeylen := 0
	for _, expr := range accumExprs {
		name, initial, val := parseKeyValInitial(expr, defaultInitial)
		if err := aggr.AddDataExpr(name, val, initial); err != nil {
			logger.Fatalf(helpers.ExitCodeInvalidUsage, "Error compiling expression %s: %s", expr, err)
		} else {
			if len(name) > maxKeylen {
				maxKeylen = len(name)
			}
		}
	}

	// Set up sorting
	var sorter = sorting.ByContextual()
	if sortReverse {
		sorter = sorting.Reverse(sorter)
	}
	if sort != "" {
		if err := aggr.SetSort(sort); err != nil {
			logger.Fatalf(helpers.ExitCodeInvalidUsage, "Error setting sort: %s", err)
		}
	}

	// run the aggregation
	if aggr.GroupColCount() > 0 || table {
		// Table output
		table := termrenderers.NewTable(vt, colCount, rowCount)

		// write header (will never shift)
		{
			rowBuf := make([]string, aggr.ColCount())
			for i, groupCol := range aggr.GroupCols() {
				rowBuf[i] = color.Wrap(color.Underline+color.BrightYellow, groupCol)
			}
			for i, dataCol := range aggr.DataCols() {
				rowBuf[aggr.GroupColCount()+i] = color.Wrap(color.Underline+color.BrightBlue, dataCol)
			}
			table.WriteRow(0, rowBuf...)
		}

		helpers.RunAggregationLoop(extractor, aggr, func() {
			// write data
			for i, group := range aggr.Groups(sorter) {
				rowBuf := make([]string, aggr.ColCount())
				data := aggr.Data(group)
				for idx, item := range group.Parts() {
					rowBuf[idx] = color.Wrap(color.BrightWhite, item)
				}
				copy(rowBuf[aggr.GroupColCount():], data)
				table.WriteRow(i+1, rowBuf...)
			}

			// write footer
			table.WriteFooter(0, helpers.FWriteExtractorSummary(extractor, aggr.ParseErrors(),
				fmt.Sprintf("(R: %d; C: %d)", aggr.DataCount(), aggr.ColCount())))
			table.WriteFooter(1, batcher.StatusString())
		})
	} else {
		// Simple output
		helpers.RunAggregationLoop(extractor, aggr, func() {
			items := aggr.Data("")
			colNames := aggr.DataCols()
			for idx, expr := range items {
				vt.WriteForLine(idx, colNames[idx]+strings.Repeat(" ", maxKeylen-len(colNames[idx]))+": "+expr)
			}
			vt.WriteForLine(len(items), helpers.FWriteExtractorSummary(extractor, aggr.ParseErrors()))
			vt.WriteForLine(len(items)+1, batcher.StatusString())
		})
	}

	vt.Close()

	if err := helpers.TryWriteCSV(c, aggr, csv.WriteAccumulator); err != nil {
		return err
	}

	return helpers.DetermineErrorState(batcher, extractor, aggr)
}

func parseKeyValInitial(s, defaultInitial string) (key, initial, val string) {
	eqSep := strings.IndexByte(s, '=')
	if eqSep < 0 {
		return s, defaultInitial, s
	}
	k := s[:eqSep]
	v := s[eqSep+1:]

	initialSep := strings.IndexByte(k, ':')
	if initialSep >= 0 {
		return k[:initialSep], k[initialSep+1:], v
	}
	return k, defaultInitial, v
}

func reduceCommand() *cli.Command {
	cmd := helpers.AdaptCommandForExtractor(cli.Command{
		Name:     "reduce",
		Action:   reduceFunction,
		Usage:    "Aggregate the results of a query based on an expression, pulling customized summary from the extracted data",
		Aliases:  []string{"r"},
		Category: cmdCatVisualize,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "accumulator",
				Aliases: []string{"a"},
				Usage:   "Specify one or more expressions to execute for each match. `{.}` is the accumulator. Syntax: `[name[:initial]=]expr`",
			},
			&cli.StringSliceFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "Specifies one or more expressions to group on. Syntax: `[name=]expr`",
			},
			&cli.StringFlag{
				Name:  "initial",
				Usage: "Specify the default initial value for any accumulators that don't specify",
				Value: "0",
			},
			&cli.BoolFlag{
				Name:  "table",
				Usage: "Force output to be a table, even when there are no groups",
			},
			&cli.IntFlag{
				Name:    "num",
				Aliases: []string{"rows", "n"},
				Usage:   "Number of elements to display",
				Value:   20,
			},
			&cli.IntFlag{
				Name:  "cols",
				Usage: "Number of columns to display",
				Value: 10,
			},
			&cli.StringFlag{
				Name:        "sort",
				Usage:       "Specify an expression to sort groups by. Will sort result in alphanumeric order",
				DefaultText: "group key",
			},
			&cli.BoolFlag{
				Name:  "sort-reverse",
				Usage: "Reverses sort order",
			},
			helpers.SnapshotFlag,
			helpers.NoOutFlag,
			helpers.CSVFlag,
		},
	})

	// Rewrite the default extraction to output array rather than {0} match
	{
		didInject := false
		for _, flag := range cmd.Flags {
			if slice, ok := flag.(*cli.StringSliceFlag); ok && slice.Name == "extract" {
				slice.Value = cli.NewStringSlice("{@}")
				didInject = true
				break
			}
		}

		if !didInject { // To catch issues in tests
			panic("Unable to inject extract change")
		}
	}

	return cmd
}
