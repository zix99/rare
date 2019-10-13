package markdowncli

import (
	"bufio"
	"fmt"
	"io"
	"rare/pkg/color"
	"regexp"
	"strings"
)

var headerColors = []color.ColorCode{color.Green, color.BrightBlue, color.Yellow, color.BrightMagenta}

var rSymbol = regexp.MustCompile("`(.*?)`")

func WriteMarkdownToTerm(out io.Writer, reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	headerDepth := 0

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") { // header
			headerDepth = strings.Count(line, "#") - 1
			headerColor := headerColors[headerDepth%len(headerColors)]
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth), color.Wrap(color.Bold, color.Wrap(headerColor, line)))
		} else if strings.Contains(line, "```") {

		} else {
			line = rSymbol.ReplaceAllStringFunc(line, func(match string) string {
				return color.Wrap(color.BrightWhite, match)
			})
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth+1), line)
		}
	}
}
