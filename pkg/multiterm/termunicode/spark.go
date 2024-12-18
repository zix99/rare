package termunicode

import (
	"io"
	"rare/pkg/multiterm/termscaler"
)

var sparkBlocks = [...]rune{
	'_',
	'\u2581',
	'\u2582',
	'\u2583',
	'\u2584',
	'\u2585',
	'\u2586',
	'\u2587',
	'\u2588',
}

var sparkAscii = [...]rune{
	'_', '.', '-', '^',
}

func SparkWrite(w io.StringWriter, scaled float64) {
	if !UnicodeEnabled {
		var blockChar = termscaler.Bucket(len(sparkAscii), scaled)
		w.WriteString(string(sparkAscii[blockChar]))
	} else {
		var blockChar = termscaler.Bucket(len(sparkBlocks), scaled)
		w.WriteString(string(sparkBlocks[blockChar]))
	}
}
