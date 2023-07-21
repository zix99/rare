package termunicode

import (
	"io"
	"rare/pkg/color"
	"rare/pkg/multiterm/termscaler"
)

const nonUnicodeBlock rune = '|'
const nonUnicodeBlockStr string = string(nonUnicodeBlock)

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

const barUnicodePartCount = len(barUnicode)

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

// Write a bar, possibly with partial runes. Not to be used with stacking
func BarWrite(w io.StringWriter, val float64, maxLen int) {
	if UnicodeEnabled {
		remainingBlocks := termscaler.LengthVal(maxLen*barUnicodePartCount, val)
		for remainingBlocks >= barUnicodePartCount {
			w.WriteString(string(fullBlock))
			remainingBlocks -= barUnicodePartCount
		}
		if remainingBlocks > 0 {
			w.WriteString(string(barUnicode[remainingBlocks]))
		}
	} else {
		blocks := termscaler.LengthVal(maxLen, val)
		for blocks > 0 {
			w.WriteString(nonUnicodeBlockStr)
			blocks--
		}
	}
}

/*
Draws various bars. Because of various outputs, there are different styles:
- Color, Unicode: Uses full/partial unicode blocks, with coloring to stack or
- NoCol, Unicode: Uses blocks if not stacked, otherwise ascii digits if stacked
- NoCol, NoUncid: Uses block-sub if not stacked, otherwise ascii digits
- Color, NoUnicd: Uses blocks with sub-char in all cases
*/

// Return the character+color to be used for a bargraph given the global context, and bar context
// useful for writing a key
func BarKey(idx int) string {
	var blockChar rune = nonUnicodeBlock
	if UnicodeEnabled {
		blockChar = fullBlock
	}

	if color.Enabled {
		return color.Wrap(color.GroupColors[idx%len(color.GroupColors)], string(blockChar))
	} else {
		return string(barAscii[idx%len(barAscii)])
	}
}

// Write a bar with a series of values, stacked with runes based on the global context
func BarWriteStacked(w io.StringWriter, maxVal, maxLen int64, vals ...int64) {
	if color.Enabled {
		// Have color, so use it as the 'key'

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
		// No color, so must use ascii char
		for i := 0; i < len(vals); i++ {
			barWriteRunes(w, barAscii[i%len(barAscii)], vals[i], maxVal, maxLen)
		}
	}
}
