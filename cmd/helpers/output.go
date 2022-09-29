package helpers

import (
	"rare/pkg/multiterm"
	"rare/pkg/multiterm/termstate"

	"github.com/urfave/cli/v2"
)

var SnapshotFlag = &cli.BoolFlag{
	Name:  "snapshot",
	Usage: "In aggregators that support it, only output final results, and not progressive updates. Will enable automatically when piping output",
}

func BuildVTerm(forceSnapshot bool) multiterm.MultilineTerm {
	if forceSnapshot || termstate.IsPipedOutput() {
		return multiterm.NewBufferedTerm()
	}
	return multiterm.New()
}

func BuildVTermFromArguments(c *cli.Context) multiterm.MultilineTerm {
	snapshot := c.Bool(SnapshotFlag.Name)
	return BuildVTerm(snapshot)
}
