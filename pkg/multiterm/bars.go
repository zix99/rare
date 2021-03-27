package multiterm

import (
	"io"
	"strings"
)

const nonUnicodeBlock rune = '|'

const fullBlock rune = '\u2588'

var UnicodeEnabled = true

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

func BarWrite(w io.StringWriter, val, maxVal, maxLen int64) {
	if val > maxVal {
		val = maxVal
	}

	if UnicodeEnabled {
		remainingBlocks := val * maxLen * barUnicodePartCount / maxVal
		for remainingBlocks >= barUnicodePartCount {
			w.WriteString(string(fullBlock))
			remainingBlocks -= barUnicodePartCount
		}
		if remainingBlocks > 0 {
			w.WriteString(string(barUnicode[remainingBlocks]))
		}
	} else {
		blocks := val * maxLen / maxVal
		for blocks > 0 {
			w.WriteString(string(nonUnicodeBlock))
			blocks--
		}
	}
}

func BarString(val, maxVal, maxLen int64) string {
	var sb strings.Builder
	BarWrite(&sb, val, maxVal, maxLen)
	return sb.String()
}
