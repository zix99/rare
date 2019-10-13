package cmd

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/urfave/cli"
)

func helpFunction(c *cli.Context) error {
	box := packr.New("Help Docs", "../docs")

	docname := c.Args().First()

	if docname == "" || docname == "list" {
		fmt.Println("Available Docs:")
		for _, name := range box.List() {
			fmt.Printf("  %s\n", strings.TrimSuffix(name, ".md"))
		}

	} else if box.Has(docname + ".md") {
		fmt.Println(box.FindString(docname + ".md"))
	} else {
		fmt.Printf("Error: No such doc %s\n", docname)
	}
	return nil
}

func HelpCommand() *cli.Command {
	return &cli.Command{
		Name:      "docs",
		Usage:     "Access help documentation",
		ArgsUsage: "[doc]",
		Action:    helpFunction,
	}
}
