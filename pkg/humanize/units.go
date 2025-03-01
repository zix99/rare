package humanize

import (
	"strconv"
)

var iecSizes = [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB"} // 1024
var siSizes = [...]string{"b", "kB", "mB", "gB", "tB", "pB", "eB", "zB"}  // 1000
var unitSize = [...]string{"", "k", "M", "B", "T"}                        // 1000

// Returns bytesize only if humanize is enabled
func ByteSize(n uint64) string {
	if !Enabled {
		return strconv.FormatUint(n, 10)
	}
	return AlwaysByteSize(n, 2)
}

// AlwaysByteSize formats bytesize (iec, power of 2) without checking `Enabled` first
func AlwaysByteSize(n uint64, precision int) string {
	return unitize(int64(n), 1024, precision, " ", iecSizes[:])
}

// Bytesize using SI (1000) units. If enabled
func ByteSizeSi(n uint64) string {
	if !Enabled {
		return strconv.FormatUint(n, 10)
	}
	return AlwaysByteSizeSi(n, 2)
}

// Bytesize using SI (1000) units, even if disabled
func AlwaysByteSizeSi(n uint64, precision int) string {
	return unitize(int64(n), 1000, precision, " ", siSizes[:])
}

// Downscale numbers by thousands (unless disabled)
func Downscale(n int64, precision int) string {
	if !Enabled {
		return strconv.FormatInt(n, 10)
	}
	return AlwaysDownscale(n, precision)
}

// Downscale number by thousands
func AlwaysDownscale(n int64, precision int) string {
	return unitize(n, 1000, precision, "", unitSize[:])
}

// downscale a number to a set of units
func unitize(n, step int64, precision int, delim string, units []string) string {
	buf := make([]byte, 0, 16)

	if n > -step && n < step {
		buf = strconv.AppendInt(buf, n, 10)
		unit := units[0]
		if len(unit) > 0 {
			buf = append(buf, delim...)
			buf = append(buf, unit...)
		}
		return string(buf)
	}

	nf, sf := float64(n), float64(step)
	rank := 0
	for (nf <= -sf || nf >= sf) && rank < len(units)-1 {
		nf /= sf
		rank++
	}

	buf = strconv.AppendFloat(buf, nf, 'f', precision, 64)
	unit := units[rank]
	if len(unit) > 0 {
		buf = append(buf, delim...)
		buf = append(buf, unit...)
	}
	return string(buf)
}
