package helpers

import (
	"testing"

	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
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

func TestModifyArgument(t *testing.T) {
	testFlag := &cli.BoolFlag{
		Name: "cat",
	}
	cmd := AdaptCommandForExtractor(cli.Command{
		Flags: []cli.Flag{
			testFlag,
		},
	})

	assert.False(t, testFlag.Value)
	ModifyArgOrPanic(cmd, "cat", func(arg *cli.BoolFlag) {
		arg.Value = true
	})
	assert.True(t, testFlag.Value)

	assert.Panics(t, func() {
		ModifyArgOrPanic(cmd, "dog", func(arg *cli.BoolFlag) {})
	})
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
	app.Commands = []*cli.Command{
		cmd,
	}

	runApp := func(args string) error {
		return app.Run(append([]string{"app", "test"}, testutil.SplitQuotedString(args)...))
	}

	assert.NoError(t, runApp(""))
	assert.NoError(t, runApp(`-I -i "{eq {0} abc}" ../testdata/log.txt`))
	assert.NoError(t, runApp(`-f ../testdata/log.txt`))
	assert.NoError(t, runApp(`-m ".*" ../testdata/log.txt`))
	assert.NoError(t, runApp(`-I -m ".*" ../testdata/log.txt`))
	assert.NoError(t, runApp(`-d "%{}" ../testdata/log.txt`))
	assert.NoError(t, runApp(`-I -d "%{}" ../testdata/log.txt`))
	testLogFatal(t, 2, func() {
		runApp("--batch 0 ../testdata/log.txt")
	})
	testLogFatal(t, 2, func() {
		runApp("--readers 0 ../testdata/log.txt")
	})
	testLogFatal(t, 2, func() {
		runApp("--readers-buffer 0 ../testdata/log.txt")
	})
	testLogFatal(t, 2, func() {
		runApp("--poll ../testdata/log.txt")
	})
	testLogFatal(t, 2, func() {
		runApp("--tail ../testdata/log.txt")
	})
	testLogFatal(t, 2, func() {
		runApp("-z -")
	})
	testLogFatal(t, 2, func() {
		runApp(`-m ".(" -`)
	})
	testLogFatal(t, 2, func() {
		runApp(`-i "{0" -`)
	})
	testLogFatal(t, 2, func() {
		runApp(`-m regex -d dissect -`)
	})
	testLogFatal(t, 2, func() {
		runApp(`-d "%{unclosed" -`)
	})
	assert.Equal(t, 7, actionCalled)
}
