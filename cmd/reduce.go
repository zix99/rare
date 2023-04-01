package cmd

import (
	"rare/cmd/helpers"
	"rare/pkg/aggregation"
	"rare/pkg/logger"
	"strings"

	"github.com/urfave/cli/v2"
)

func reduceFunction(c *cli.Context) error {
	var (
		accumExprs = c.StringSlice("accumulator")
		initial    = c.String("initial")
		extra      = c.Bool("extra")
	)

	vt := helpers.BuildVTermFromArguments(c)
	batcher := helpers.BuildBatcherFromArguments(c)
	extractor := helpers.BuildExtractorFromArguments(c, batcher)

	aggr := aggregation.NewExprAccumulatorSet()
	maxKeylen := 0
	for _, expr := range accumExprs {
		name, val := parseKeyValue(expr)
		if err := aggr.Add(name, val, initial); err != nil {
			logger.Printf("Error compiling expression %s: %s", expr, err)
		} else {
			if len(name) > maxKeylen {
				maxKeylen = len(name)
			}
		}
	}

	helpers.RunAggregationLoop(extractor, aggr, func() {
		items := aggr.Items()
		for idx, expr := range items {
			if extra {
				vt.WriteForLine(idx, expr.Name+strings.Repeat(" ", maxKeylen-len(expr.Name))+": "+expr.Accum.Value())
			} else {
				vt.WriteForLine(idx, expr.Accum.Value())
			}
		}
		vt.WriteForLine(len(items), helpers.FWriteExtractorSummary(extractor, aggr.ParseErrors()))
		vt.WriteForLine(len(items)+1, batcher.StatusString())
	})

	vt.Close()

	return helpers.DetermineErrorState(batcher, extractor, aggr)
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
				Usage:   "Specify one or more expressions to execute for each match. `{.}` is the accumulator",
			},
			&cli.BoolFlag{
				Name:    "extra",
				Aliases: []string{"x"},
				Usage:   "Write out the result keys as well as their values",
			},
			&cli.StringFlag{
				Name:  "initial",
				Usage: "Specify the initial value for any accumulators",
				Value: "0",
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
