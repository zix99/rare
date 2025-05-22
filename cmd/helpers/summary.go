package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"rare/pkg/color"
	"rare/pkg/extractor"
	"rare/pkg/extractor/batchers"
	"rare/pkg/extractor/dirwalk"
	"rare/pkg/humanize"
)

func FWriteMatchSummary(w io.Writer, matched, total uint64) {
	fmt.Fprintf(w, "Matched: %s / %s",
		color.Wrap(color.BrightGreen, humanize.Hui(matched)),
		color.Wrap(color.BrightWhite, humanize.Hui(total)))
}

func FWriteExtractorSummary(extractor *extractor.Extractor, errors uint64, additionalParts ...string) string {
	var w bytes.Buffer
	FWriteMatchSummary(&w, extractor.MatchedLines(), extractor.ReadLines())
	for _, p := range additionalParts {
		w.WriteRune(' ')
		w.WriteString(p)
	}
	if extractor.IgnoredLines() > 0 {
		fmt.Fprintf(&w, " (Ignored: %s)", color.Wrap(color.Red, humanize.Hui(extractor.IgnoredLines())))
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

func WriteBatcherSummary(w io.Writer, b *batchers.Batcher, walker dirwalk.Metrics) {
	fmt.Fprintf(w, "Read   : %s file(s) (%s)",
		color.Wrap(color.BrightWhite, humanize.Hi32(b.ReadFiles())),
		color.Wrap(color.BrightBlue, humanize.ByteSize(b.ReadBytes())),
	)

	if walker != nil {
		if skipped := walker.ExcludedCount(); skipped > 0 {
			fmt.Fprintf(w, ", %s excluded", color.Wrap(color.Yellow, humanize.Hui(skipped)))
		}
	}
	if errCount := b.ReadErrors(); errCount > 0 {
		fmt.Fprintf(w, ", %s error(s)", color.Wrap(color.Red, humanize.Hi32(errCount)))
	}
	fmt.Fprint(w, "\n")
}
