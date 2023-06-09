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
		&cli.BoolFlag{
			Name: "snapshot",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		BuildVTermFromArguments(ctx)
		return nil
	}
	assert.NoError(t, app.Run([]string{"", "--snapshot"}))
}

func TestTryWriteCSV(t *testing.T) {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "csv",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		agg := aggregation.NewCounter()
		TryWriteCSV(ctx, agg, csv.WriteCounter)
		return nil
	}
	assert.NoError(t, app.Run([]string{""}))
	assert.NoError(t, app.Run([]string{"", "--csv", "-"}))
}
