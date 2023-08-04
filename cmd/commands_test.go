package cmd

import (
	"fmt"
	"os"
	"rare/pkg/logger"
	"rare/pkg/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestCommandGetter(t *testing.T) {
	assert.NotNil(t, GetSupportedCommands())
}

func testCommandSet(t *testing.T, command *cli.Command, commandArgList ...string) {
	for _, args := range commandArgList {
		assert.NoError(t, testCommand(command, args))
	}
}

func testCommandCapture(command *cli.Command, cmd string) (stdout, stderr string, err error) {
	return testutil.Capture(func(w *os.File) error {
		return testCommand(command, cmd)
	})
}

func testCommand(command *cli.Command, cmd string) error {
	app := cli.NewApp()

	command.Name = "_testcommand"
	app.Commands = []*cli.Command{
		command,
	}
	app.ExitErrHandler = func(context *cli.Context, err error) {
		// disabled failure
		fmt.Fprint(os.Stderr, err.Error())
	}

	commandArgs := append([]string{"app", "_testcommand"}, testutil.SplitQuotedString(cmd)...)

	return app.Run(commandArgs)
}

// Cause logger.fatal* to result in panic() for testability
func catchLogFatal(t *testing.T, expectsCode int, f func()) (code int) {
	code = -1

	oldExit := logger.OsExit
	defer func() {
		logger.OsExit = oldExit
	}()
	logger.OsExit = func(v int) {
		code = v
		panic("logger.osexit")
	}

	assert.PanicsWithValue(t, "logger.osexit", f)
	assert.Equal(t, expectsCode, code)
	return
}
