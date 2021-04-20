package cmd

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestCommandGetter(t *testing.T) {
	assert.NotNil(t, GetSupportedCommands())
}

func testCommandSet(t *testing.T, command *cli.Command, commandArgList ...string) {
	for _, args := range commandArgList {
		assert.NoError(t, testCommand(command, args))
	}
}

func testCommand(command *cli.Command, cmd string) error {
	app := cli.NewApp()

	command.Name = "_testcommand"
	app.Commands = []cli.Command{
		*command,
	}

	commandArgs := append([]string{"app", "_testcommand"}, splitQuotedString(cmd)...)

	return app.Run(commandArgs)
}

var stringSplitter = regexp.MustCompile(`([^\s"]+)|"([^"]*)"`)

func splitQuotedString(s string) []string {
	matches := stringSplitter.FindAllStringSubmatch(s, -1)

	ret := make([]string, 0)
	for _, v := range matches {
		if v[2] != "" {
			ret = append(ret, v[2])
		} else {
			ret = append(ret, v[1])
		}
	}

	return ret
}
