package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestBuildVTerm(t *testing.T) {
	assert.NotNil(t, BuildOutVTerm(false))
	assert.NotNil(t, BuildOutVTerm(true))
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
	app.Run([]string{"", "--snapshot"})
}
