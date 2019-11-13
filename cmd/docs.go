package cmd

import (
	"fmt"
	"os"
	"rare/pkg/markdowncli"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/urfave/cli"
)

func docsFunction(c *cli.Context) error {
	box := packr.New("Help Docs", "../docs")

	docname := strings.ToLower(c.Args().First())

	if docname == "" || docname == "list" {
		fmt.Println("Available Docs:")
		for _, name := range box.List() {
			fmt.Printf("  %s\n", strings.Title(strings.TrimSuffix(name, ".md")))
		}

	} else if box.Has(docname + ".md") {
		file, err := box.Resolve(docname + ".md")
		if err != nil {
			fmt.Println(err)
		} else {
			markdowncli.WriteMarkdownToTerm(os.Stdout, file)
		}
	} else {
		fmt.Printf("Error: No such doc %s\n", docname)
	}
	return nil
}

func docsCommand() *cli.Command {
	return &cli.Command{
		Name:      "docs",
		Usage:     "Access help documentation",
		ArgsUsage: "[doc]",
		Action:    docsFunction,
	}
}
