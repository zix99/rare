package helpers

import (
	"github.com/zix99/rare/pkg/csv"
	"github.com/zix99/rare/pkg/multiterm"
	"github.com/zix99/rare/pkg/multiterm/termstate"

	"github.com/urfave/cli/v2"
)

var SnapshotFlag = &cli.BoolFlag{
	Name:     "snapshot",
	Usage:    "In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output",
	Category: cliCategoryOutput,
}

var CSVFlag = &cli.StringFlag{
	Name:     "csv",
	Aliases:  []string{"o"},
	Usage:    "Write final results to csv. Use - to output to stdout",
	Category: cliCategoryOutput,
}

var NoOutFlag = &cli.BoolFlag{
	Name:     "noout",
	Usage:    "Don't output any aggregation to stdout",
	Category: cliCategoryOutput,
}

func BuildVTerm(forceSnapshot bool) multiterm.MultilineTerm {
	if forceSnapshot || termstate.IsPipedOutput() {
		return multiterm.NewBufferedTerm()
	}
	return multiterm.New()
}

func BuildVTermFromArguments(c *cli.Context) multiterm.MultilineTerm {
	if c.Bool(NoOutFlag.Name) || c.String(CSVFlag.Name) == "-" {
		return &multiterm.NullTerm{}
	}

	snapshot := c.Bool(SnapshotFlag.Name)
	return BuildVTerm(snapshot)
}

func TryWriteCSV[T any](c *cli.Context, agg T, writer func(w csv.CSV, agg T) error) error {
	if filename := c.String(CSVFlag.Name); filename != "" {
		if w, err := csv.OpenCSV(filename); err != nil {
			return cli.Exit(err, ExitCodeOutputError)
		} else {
			defer w.Close()
			if err := writer(w, agg); err != nil {
				return cli.Exit(err, ExitCodeOutputError)
			}
		}
	}
	return nil
}
