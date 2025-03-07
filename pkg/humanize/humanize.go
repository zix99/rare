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
