package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestAdaptingCommandForExtractor(t *testing.T) {
	called := false
	dummyAfter := func(c *cli.Context) error {
		called = true
		return nil
	}
	cmd := cli.Command{
		Name:        "test",
		Description: "this is a test",
		After:       dummyAfter,
	}

	adapted := AdaptCommandForExtractor(cmd)
	assert.NotNil(t, adapted.After)
	assert.NotNil(t, adapted.ArgsUsage)
	assert.Equal(t, "test", adapted.Name)

	adapted.After(nil)
	assert.True(t, called)
}

func TestBuildingExtractorFromContext(t *testing.T) {
	actionCalled := 0
	cmdAction := func(c *cli.Context) error {
		batcher := BuildBatcherFromArguments(c)
		extractor := BuildExtractorFromArguments(c, batcher)
		assert.NotNil(t, extractor)

		actionCalled++
		return nil
	}
	cmd := AdaptCommandForExtractor(cli.Command{
		Name:   "test",
		Action: cmdAction,
	})

	app := cli.NewApp()
	app.Commands = []cli.Command{
		*cmd,
	}

	app.Run([]string{"app", "test"})
	app.Run([]string{"app", "test", "-i", "{eq {0} abc}", "../testdata/log.txt"})
	app.Run([]string{"app", "test", "-f", "../testdata/log.txt"})
	assert.Equal(t, 3, actionCalled)
}
