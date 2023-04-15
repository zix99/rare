package termunicode

import (
	"io"
	"rare/pkg/color"
	"strings"
)

const nonUnicodeBlock rune = '|'

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

var barAscii = [...]rune{
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
}

const barUnicodePartCount int64 = int64(len(barUnicode))

// write a length of runes for a given bar parameters
func barWriteRunes(w io.StringWriter, blockChar rune, val, maxVal, maxLen int64) {
	if val > maxVal {
		val = maxVal
	}

	blocks := val * maxLen / maxVal
	for blocks > 0 {
		w.WriteString(string(blockChar))
		blocks--
	}
}

// Return the character to be used for a bargraph given the global context, and bar context
// useful for writing a key
func BarKeyChar(stacked bool, idx int) string {
	if color.Enabled {
		var blockChar rune = nonUnicodeBlock
		if UnicodeEnabled {
			blockChar = fullBlock
		}
		return color.Wrap(color.GroupColors[idx%len(color.GroupColors)], string(blockChar))
	} else {
		if stacked {
			return string(barAscii[idx%len(barAscii)])
		}
		return string(fullBlock)
	}
}

// BarWriteFull does not write partial bars to the end. Useful for stacking
func BarWriteFull(w io.StringWriter, val, maxVal, maxLen int64) {
	var blockChar rune = nonUnicodeBlock
	if UnicodeEnabled {
		blockChar = fullBlock
	}

	barWriteRunes(w, blockChar, val, maxVal, maxLen)
}

// Write a bar, possibly with partial runes. Not to be used with stacking
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

// Write a bar with a series of values, stacked with runes based on the global context
func BarWriteStacked(w io.StringWriter, maxVal, maxLen int64, vals ...int64) {
	if color.Enabled {
		var blockChar rune = nonUnicodeBlock
		if UnicodeEnabled {
			blockChar = fullBlock
		}

		for i := 0; i < len(vals); i++ {
			color.Write(w, color.GroupColors[i%len(color.GroupColors)], func(w io.StringWriter) {
				barWriteRunes(w, blockChar, vals[i], maxVal, maxLen)
			})
		}
	} else {
		for i := 0; i < len(vals); i++ {
			barWriteRunes(w, barAscii[i%len(barAscii)], vals[i], maxVal, maxLen)
		}
	}
}

// BarWrite, but to a string
func BarString(val, maxVal, maxLen int64) string {
	var sb strings.Builder
	BarWrite(&sb, val, maxVal, maxLen)
	return sb.String()
}
