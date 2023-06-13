package helpers

import (
	"rare/pkg/csv"
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termstate"

	"github.com/urfave/cli/v2"
)

var SnapshotFlag = &cli.BoolFlag{
	Name:  "snapshot",
	Usage: "In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output",
}

var CSVFlag = &cli.StringFlag{
	Name:    "csv",
	Aliases: []string{"o"},
	Usage:   "Write final results to csv. Use - to output to stdout",
}

var NoOutFlag = &cli.BoolFlag{
	Name:  "noout",
	Usage: "Don't output any aggregation to stdout",
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
			return err
		} else {
			defer w.Close()
			return writer(w, agg)
		}
	}
	return nil
}
