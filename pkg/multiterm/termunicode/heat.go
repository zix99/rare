package termunicode

import (
	"io"
	"rare/pkg/color"
)

const heatmapEscape = "\x1b[38;5;"

var heatmapColors = []color.ColorCode{
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

var heatmapNonUnicode rune = '#'

func HeatWriteLinear(w io.StringWriter, val, min, max int64) {
	if val >= max {
		val = max - 1
	}
	if val < min {
		val = min
	}

	var blockChar = heatmapNonUnicode
	if UnicodeEnabled {
		blockChar = fullBlock
	}

	if max-min <= 0 {
		w.WriteString(color.Wrap(heatmapColors[0], string(blockChar)))
	} else {
		blockIdx := ((val - min) * int64(len(heatmapColors))) / (max - min)
		hc := heatmapColors[blockIdx]

		w.WriteString(color.Wrap(hc, string(blockChar)))
	}
}
