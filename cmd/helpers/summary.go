package helpers

import (
	"bytes"
	"fmt"
	"os"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
)

func FWriteExtractorSummary(extractor *extractor.Extractor, errors uint64, additionalParts ...string) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "Matched: %s / %s",
		color.Wrapi(color.BrightGreen, humanize.Hi(extractor.MatchedLines())),
		color.Wrapi(color.BrightWhite, humanize.Hi(extractor.ReadLines())))
	for _, p := range additionalParts {
		w.WriteString(p)
	}
	if extractor.IgnoredLines() > 0 {
		fmt.Fprintf(&w, " (Ignored: %s)", color.Wrapi(color.Red, humanize.Hi(extractor.IgnoredLines())))
	}
	if errors > 0 {
		fmt.Fprintf(&w, " %s", color.Wrapf(color.Red, "(Errors: %v)", humanize.Hi(errors)))
	}
	return w.String()
}

func WriteExtractorSummary(extractor *extractor.Extractor) {
	os.Stderr.WriteString(FWriteExtractorSummary(extractor, 0))
	os.Stderr.WriteString("\n")
}
