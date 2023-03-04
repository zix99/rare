package humanize

import (
	"strconv"
)

// Enabled determines whether to use language message printer, or fmt
var Enabled = true
var Decimals = 4

func Hi(arg int64) string {
	if !Enabled {
		return strconv.FormatInt(arg, 10)
	}
	return humanizeInt(arg)
}

func Hui(arg uint64) string {
	if !Enabled {
		return strconv.FormatUint(arg, 10)
	}
	return humanizeInt(arg)
}

func Hi32(arg int) string {
	if !Enabled {
		return strconv.Itoa(arg)
	}
	return humanizeInt(arg)
}

func Hf(arg float64) string {
	return Hfd(arg, Decimals)
}

func Hfd(arg float64, decimals int) string {
	if !Enabled {
		return strconv.FormatFloat(arg, 'f', decimals, 64)
	}
	return humanizeFloat(arg, decimals)
}

var byteSizes = [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB"}

func ByteSize(n uint64) string {
	if !Enabled {
		return strconv.FormatUint(n, 10)
	}
	return AlwaysByteSize(n, 2)
}

// AlwaysByteSize formats bytesize without checking `Enabled` first
func AlwaysByteSize(n uint64, precision int) string {
	if n < 1024 { // Never a decimal for byte-unit
		return strconv.FormatUint(n, 10) + " " + byteSizes[0]
	}

	var nf float64 = float64(n)
	labelIdx := 0
	for nf >= 1024.0 && labelIdx < len(byteSizes)-1 {
		nf /= 1024.0
		labelIdx++
	}

	return strconv.FormatFloat(nf, 'f', precision, 64) + " " + byteSizes[labelIdx]
}
