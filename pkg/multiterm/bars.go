package multiterm

import (
	"io"
)

const fullBlock rune = '\u2588'

var barUnicode = [...]rune{
	'\u0000',
	'\u258f',
	'\u258e',
	'\u258d',
	'\u258c',
	'\u258b',
	'\u258a',
	'\u2589',
	'\u2588',
}
var barUnicodePartCount int64 = int64(len(barUnicode))

type RuneWriter interface {
	io.Writer
	WriteRune(r rune) (int, error)
}

func WriteBar(w io.StringWriter, val, maxVal, maxLen int64) {
	if val > maxVal {
		val = maxVal
	}

	remainingBlocks := val * maxLen * barUnicodePartCount / maxVal
	for remainingBlocks >= barUnicodePartCount {
		w.WriteString(string(fullBlock))
		remainingBlocks -= barUnicodePartCount
	}
	if remainingBlocks > 0 {
		w.WriteString(string(barUnicode[remainingBlocks]))
	}
}
