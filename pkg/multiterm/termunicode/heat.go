package termunicode

import (
	"io"
	"rare/pkg/color"
)

const heatmapEscape = "\x1b[38;5;"

var heatmapColors = [...]color.ColorCode{
	heatmapEscape + "16m",
	heatmapEscape + "17m",
	heatmapEscape + "18m",
	heatmapEscape + "19m",
	heatmapEscape + "20m",
	heatmapEscape + "21m",
	heatmapEscape + "57m",
	heatmapEscape + "93m",
	heatmapEscape + "129m",
	heatmapEscape + "165m",
	heatmapEscape + "201m",
	heatmapEscape + "200m",
	heatmapEscape + "199m",
	heatmapEscape + "198m",
	heatmapEscape + "197m",
	heatmapEscape + "196m",
}

const heatmapColorsLen int64 = int64(len(heatmapColors))

var heatmapAscii = [...]string{
	"-",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
}

const heatmapAsciiLen int64 = int64(len(heatmapAscii))

const heatmapNonUnicode rune = '#'

func HeatWriteLinear(w io.StringWriter, val, min, max int64) {
	if val > max {
		val = max
	}
	if val < min {
		val = min
	}

	if !color.Enabled {
		// Fallback to numeric single-digit display when no colors are available
		if max <= min {
			w.WriteString(heatmapAscii[0])
		} else {
			idx := ((val - min) * (heatmapAsciiLen - 1)) / (max - min)
			w.WriteString(heatmapAscii[idx])
		}
	} else {
		var blockChar = heatmapNonUnicode
		if UnicodeEnabled {
			blockChar = fullBlock
		}

		if max <= min {
			w.WriteString(color.Wrap(heatmapColors[0], string(blockChar)))
		} else {
			blockIdx := ((val - min) * (heatmapColorsLen - 1)) / (max - min)
			hc := heatmapColors[blockIdx]

			w.WriteString(color.Wrap(hc, string(blockChar)))
		}
	}
}
