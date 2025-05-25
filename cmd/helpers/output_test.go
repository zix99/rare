package helpers

import (
	"rare/pkg/aggregation"
	"rare/pkg/csv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestBuildVTerm(t *testing.T) {
	assert.NotNil(t, BuildVTerm(false))
	assert.NotNil(t, BuildVTerm(true))
}

func TestBuildVTermFromArgs(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		SnapshotFlag,
		NoOutFlag,
	}
	app.Action = func(ctx *cli.Context) error {
		BuildVTermFromArguments(ctx)
		return nil
	}
	assert.NoError(t, app.Run([]string{"", "--snapshot"}))
	assert.NoError(t, app.Run([]string{"", "--noout"}))
}

func TestTryWriteCSV(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		CSVFlag,
	}
	app.Action = func(ctx *cli.Context) error {
		agg := aggregation.NewCounter()
		return TryWriteCSV(ctx, agg, csv.WriteCounter)
	}
	app.ExitErrHandler = func(cCtx *cli.Context, err error) {}
	assert.NoError(t, app.Run([]string{""}))
	assert.NoError(t, app.Run([]string{"", "--csv", "-"}))
	assert.Error(t, app.Run([]string{"", "--csv", "/!@#bad-filename"}))
}
