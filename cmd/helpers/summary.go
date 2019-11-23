package helpers

import (
	"bytes"
	"fmt"
	"os"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
)

func FWriteExtractorSummary(extractor *extractor.Extractor, additionalFormat string, args ...interface{}) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "Matched: %s / %s",
		color.Wrapi(color.BrightGreen, humanize.Hi(extractor.MatchedLines())),
		color.Wrapi(color.BrightWhite, humanize.Hi(extractor.ReadLines())))
	if additionalFormat != "" {
		fmt.Fprintf(&w, additionalFormat, args...)
	}
	if extractor.IgnoredLines() > 0 {
		fmt.Fprintf(&w, " (Ignored: %s)", color.Wrapi(color.Red, humanize.Hi(extractor.IgnoredLines())))
	}
	return w.String()
}

func WriteExtractorSummary(extractor *extractor.Extractor) {
	os.Stderr.WriteString(FWriteExtractorSummary(extractor, ""))
	os.Stderr.WriteString("\n")
}
