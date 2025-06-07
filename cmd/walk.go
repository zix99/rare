package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/zix99/rare/cmd/helpers"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/humanize"

	"github.com/urfave/cli/v2"
)

func walkFunction(c *cli.Context) error {
	walk := helpers.BuildPathWalkerFromArguments(c)
	paths := c.Args().Slice()

	stdout := bufio.NewWriter(os.Stdout)

	for path := range walk.Walk(paths...) {
		stdout.WriteString(path)
		stdout.WriteRune('\n')
	}
	stdout.Flush()

	if count := walk.TotalCount(); count > 0 {
		fmt.Fprintf(os.Stderr, "Found %s path(s)", color.Wrap(color.BrightGreen, humanize.Hui(count)))
		if excluded := walk.ExcludedCount(); excluded > 0 {
			fmt.Fprintf(os.Stderr, ", %s excluded", color.Wrap(color.Yellow, humanize.Hui(excluded)))
		}
		if errors := walk.ErrorCount(); errors > 0 {
			fmt.Fprintf(os.Stderr, ", %s error(s)", color.Wrap(color.Red, humanize.Hui(errors)))
		}
		fmt.Fprint(os.Stderr, "\n")
	} else {
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
