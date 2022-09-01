package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestAddSortingCommands(t *testing.T) {
	cmd := &cli.Command{}
	assert.Panics(t, func() {
		AddSortFlag(cmd, "no-exist")
	})

	AddSortFlag(cmd, "value")

	assert.Equal(t, "sort", cmd.Flags[0].Names()[0])
	assert.Equal(t, "sort-reverse", cmd.Flags[1].Names()[0])
}

func TestBuildSorter(t *testing.T) {
	assert.NotNil(t, BuildSorter("text", false))
	assert.NotNil(t, BuildSorter("smart", false))
	assert.NotNil(t, BuildSorter("contextual", false))
	assert.NotNil(t, BuildSorter("value", false))
	assert.NotNil(t, BuildSorter("value", true))
}

func TestBuildSorterFromFlags(t *testing.T) {
	app := cli.NewApp()
	cmd := &cli.Command{
		Name: "test",
		Action: func(ctx *cli.Context) error {
			sorter := BuildSorterFromFlags(ctx)
			assert.NotNil(t, sorter)
			return nil
		},
	}
	AddSortFlag(cmd, "text")
	app.Commands = []*cli.Command{cmd}

	assert.NoError(t, app.Run([]string{"", "test", "--sort=smart"}))
}
