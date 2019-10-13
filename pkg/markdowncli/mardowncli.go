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

const (
	tokenCode   = "```"
	tokenHeader = "#"
)

// WriteMarkdownToTerm does pseudo-markdown formatting
//   it doesn't follow correctly to the spec, but is close enough for our docs
func WriteMarkdownToTerm(out io.Writer, reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	headerDepth := 0
	isCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, tokenHeader) && !isCodeBlock { // header
			headerDepth = strings.Count(line, tokenHeader) - 1
			headerColor := headerColors[headerDepth%len(headerColors)]
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth), color.Wrap(color.Bold, color.Wrap(headerColor, line)))
		} else if strings.HasPrefix(line, tokenCode) {
			codeType := strings.TrimPrefix(line, tokenCode)
			isCodeBlock = !isCodeBlock

			if isCodeBlock && codeType != "" {
				fmt.Fprint(out, strings.Repeat(" ", headerDepth+1))
				fmt.Fprint(out, color.Wrapf(color.Underline, "Code %s:\n", codeType))
			}

			if isCodeBlock {
				headerDepth++
			} else {
				headerDepth--
			}
		} else {
			if isCodeBlock {
				line = color.Wrap(color.BrightMagenta, line)
			}
			line = rSymbol.ReplaceAllStringFunc(line, func(match string) string {
				return color.Wrap(color.BrightWhite, match)
			})
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth+1), line)
		}
	}
}
