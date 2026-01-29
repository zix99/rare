package markdowncli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/zix99/rare/pkg/color"
)

var headerColors = []color.ColorCode{color.Green, color.BrightBlue, color.Yellow, color.BrightMagenta}

const (
	tokenCode   = "```"
	tokenHeader = "#"
	tokenNote   = "!!!"
)

// WriteMarkdownToTerm does pseudo-markdown formatting
//
//	it doesn't follow correctly to the spec, but is close enough for our docs
func WriteMarkdownToBuf(out io.Writer, reader io.Reader) {
	scanner := bufio.NewScanner(reader)

	headerDepth := 0
	isCodeBlock := false
	inNoteBlock := false
	isFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" && headerDepth == 0 { // skip frontmatter
			isFrontmatter = !isFrontmatter
		} else if isFrontmatter {
			continue
		} else if strings.HasPrefix(line, tokenHeader) && !isCodeBlock && !inNoteBlock { // header
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
		} else if strings.HasPrefix(line, tokenNote) { // note block begin
			inNoteBlock = true
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth+1), color.Wrap(color.BrightCyan, line))
		} else if inNoteBlock && line == "" { // note block end
			inNoteBlock = false
			fmt.Fprint(out, "\n")
		} else {
			if inNoteBlock {
				line = color.Wrap(color.BrightBlack, line)
			} else if isCodeBlock {
				line = color.Wrap(color.BrightMagenta, line)
			} else {
				for _, replacer := range regexReplacement {
					line = replacer.match.ReplaceAllStringFunc(line, replacer.process)
				}
			}
			fmt.Fprintf(out, "%s%s\n", strings.Repeat(" ", headerDepth+1), line)
		}
	}
}
