package cmd

import (
	"bufio"
	"fmt"
	"os"
	"rare/cmd/helpers"
	"rare/pkg/color"
	"rare/pkg/humanize"

	"github.com/urfave/cli/v2"
)

func walkFunction(c *cli.Context) error {
	walk := helpers.BuildPathWalkerFromArguments(c)
	paths := c.Args().Slice()

	stdout := bufio.NewWriter(os.Stdout)

	var count uint64
	for path := range walk.Walk(paths...) {
		stdout.WriteString(path)
		stdout.WriteRune('\n')
		count++
	}
	stdout.Flush()

	fmt.Fprintf(os.Stderr, "Found %s path(s)\n", color.Wrap(color.BrightGreen, humanize.Hui(count)))

	if count == 0 {
		return cli.Exit("No paths found", helpers.ExitCodeNoData)
	}

	return nil
}

func walkCommand() *cli.Command {
	return &cli.Command{
		Name:        "walk",
		Usage:       "Output paths discovered via traverse rules",
		Description: "Find files using include/exclude and traversal rules by outputting paths visited",
		ArgsUsage:   "<paths...>",
		Action:      walkFunction,
		Category:    cmdCatHelp,
		Flags:       helpers.GetWalkerFlags(),
	}
}
