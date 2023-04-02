package cmd

import (
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/aggregation/sorting"
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
	)

	vt := helpers.BuildVTermFromArguments(c)
	batcher := helpers.BuildBatcherFromArguments(c)
	extractor := helpers.BuildExtractorFromArguments(c, batcher)

	aggr := aggregation.NewAccumulatingGroup(stdlib.NewStdKeyBuilder())

	// Set up groups
	for _, group := range groupExpr {
		name, val := parseKeyValue(group)
		if err := aggr.AddGroupExpr(name, val); err != nil {
			logger.Fatalf("Error compiling group expression %s: %s", group, err)
		}
	}

	// Set up expressions
	maxKeylen := 0
	for _, expr := range accumExprs {
		name, initial, val := parseKeyValInitial(expr, defaultInitial)
		if err := aggr.AddDataExpr(name, val, initial); err != nil {
			logger.Fatalf("Error compiling expression %s: %s", expr, err)
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
			logger.Fatalf("Error setting sort: %s", err)
		}
	}

	// run the aggregation
	if aggr.GroupColCount() > 0 || table {
		table := termrenderers.NewTable(vt, 10, 10) // TODO: Undo hardcode

		helpers.RunAggregationLoop(extractor, aggr, func() {
			cols := append(aggr.GroupCols(), aggr.DataCols()...)
			table.WriteRow(0, cols...)
			for i, group := range aggr.Groups(sorter) {
				data := aggr.Data(group)
				table.WriteRow(i+1, append(group.Parts(), data...)...)
			}

			table.WriteFooter(0, helpers.FWriteExtractorSummary(extractor, aggr.ParseErrors()))
			table.WriteFooter(1, batcher.StatusString())
		})
	} else {
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
		Category: cmdCatAnalyze,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:    "accumulator",
				Aliases: []string{"a"},
				Usage:   "Specify one or more expressions to execute for each match. `{.}` is the accumulator. `[name[:initial]=]expr`",
			},
			&cli.StringSliceFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "Specifies one or more expressions to group on",
			},
			&cli.StringFlag{
				Name:  "initial",
				Usage: "Specify the default initial value for any accumulators that don't specify",
				Value: "0",
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
		},
	})

	// Rewrite the default extraction to output array rather than {0} match
	for _, flag := range cmd.Flags {
		if slice, ok := flag.(*cli.StringSliceFlag); ok && slice.Name == "extract" {
			slice.Value = cli.NewStringSlice("{@}")
			break
		}
	}

	return cmd
}
