package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
)

func FWriteMatchSummary(w io.Writer, matched, total uint64) {
	fmt.Fprintf(w, "Matched: %s / %s",
		color.Wrapi(color.BrightGreen, humanize.Hui(matched)),
		color.Wrapi(color.BrightWhite, humanize.Hui(total)))
}

func FWriteExtractorSummary(extractor *extractor.Extractor, errors uint64, additionalParts ...string) string {
	var w bytes.Buffer
	FWriteMatchSummary(&w, extractor.MatchedLines(), extractor.ReadLines())
	for _, p := range additionalParts {
		w.WriteRune(' ')
		w.WriteString(p)
	}
	if extractor.IgnoredLines() > 0 {
		fmt.Fprintf(&w, " (Ignored: %s)", color.Wrapi(color.Red, humanize.Hui(extractor.IgnoredLines())))
	}
	if errors > 0 {
		fmt.Fprintf(&w, " %s", color.Wrapf(color.Red, "(Errors: %v)", humanize.Hui(errors)))
	}
	return w.String()
}

func WriteExtractorSummary(extractor *extractor.Extractor) {
	os.Stderr.WriteString(FWriteExtractorSummary(extractor, 0))
	os.Stderr.WriteString("\n")
}
