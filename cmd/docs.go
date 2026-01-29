package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strings"

	"github.com/zix99/rare/cmd/helpers"
	"github.com/zix99/rare/docs"
	"github.com/zix99/rare/pkg/color"
	"github.com/zix99/rare/pkg/markdowncli"

	"github.com/urfave/cli/v2"
)

func docsFunction(c *cli.Context) error {
	docname := strings.ToLower(c.Args().First())

	if docname == "" || docname == "list" {
		listDocFiles()
	} else if file, err := openDocFileByPartialName(docname); err == nil {
		defer file.Close()

		var buf bytes.Buffer
		markdowncli.WriteMarkdownToBuf(&buf, file)
		if c.Bool("no-pager") || helpers.TryWritePager(&buf) != nil {
			io.Copy(os.Stdout, &buf)
		}
	} else {
		return cli.Exit(fmt.Sprintf("No such doc '%s'", docname), helpers.ExitCodeInvalidUsage)
	}

	return nil
}

func listDocFiles() {
	fmt.Println(color.Wrap(color.Bold, "Available Docs:"))

	type docInfo struct {
		name         string
		summary      string
		order, depth int
	}

	entries, _ := docs.DocFS.ReadDir(docs.BasePath)
	docList := make([]docInfo, 0, len(entries))
	maxNameLen := 1
	for _, entry := range entries {
		info := docInfo{
			name: strings.TrimSuffix(entry.Name(), ".md"),
		}
		maxNameLen = max(maxNameLen, len(info.name))

		r, err := docs.DocFS.Open(docs.BasePath + "/" + entry.Name())
		if err == nil {
			defer r.Close()
			frontmatter := markdowncli.ExtractFrontmatter(r)
			info.summary = frontmatter.Description()
			info.order = frontmatter.Order()
			info.depth = frontmatter.Depth()
		}

		docList = append(docList, info)
	}

	sort.Slice(docList, func(i, j int) bool {
		di, dj := docList[i], docList[j]
		if di.order != dj.order {
			return di.order < dj.order
		}
		if di.depth != dj.depth {
			return di.depth < dj.depth
		}
		return di.name < dj.name
	})

	for _, d := range docList {
		fmt.Print(strings.Repeat("  ", d.depth+1))
		fmt.Printf("%s%s", color.Wrap(color.BrightWhite, d.name), strings.Repeat(" ", maxNameLen-len(d.name)))
		if d.summary != "" {
			fmt.Print("  ", d.summary)
		}
		fmt.Println()
	}
}

func openDocFileByPartialName(s string) (fs.File, error) {
	s = strings.ToLower(s)

	// Try exact match
	if f, err := docs.DocFS.Open(docs.BasePath + "/" + s + ".md"); err == nil {
		return f, nil
	}

	// Search for prefix of name
	entries, _ := docs.DocFS.ReadDir(docs.BasePath)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), s) {
			fullPath := docs.BasePath + "/" + entry.Name()
			f, err := docs.DocFS.Open(fullPath)
			if err == nil {
				return f, nil
			}
		}
	}

	return nil, errors.New("no such file")
}

func docsCommand() *cli.Command {
	return &cli.Command{
		Name:      "docs",
		Usage:     "Access detailed documentation",
		ArgsUsage: "[doc]",
		Action:    docsFunction,
		Category:  cmdCatHelp,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "no-pager",
				Aliases: []string{"n"},
				Usage:   "Don't use pager to view documentation",
			},
		},
	}
}
