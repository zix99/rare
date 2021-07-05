package humanize

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer = message.NewPrinter(language.English)

// Enabled determines whether to use language message printer, or fmt
var Enabled = true
var Decimals = 4

// H humanizes the output
func H(format string, args ...interface{}) string {
	if !Enabled {
		return fmt.Sprintf(format, args...)
	}
	return printer.Sprintf(format, args...)
}

func Hi(arg interface{}) string {
	if !Enabled {
		return fmt.Sprintf("%d", arg)
	}
	return printer.Sprintf("%d", arg)
}

func Hf(arg interface{}) string {
	if !Enabled {
		return fmt.Sprintf("%f", arg)
	}
	return printer.Sprintf("%.[2]*[1]f", arg, Decimals)
}

func Hfd(arg interface{}, decimals int) string {
	if !Enabled {
		return fmt.Sprintf("%f", arg)
	}
	return printer.Sprintf("%.[2]*[1]f", arg, decimals)
}

var byteSizes = [...]string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB"}

func ByteSize(n uint64) string {
	if !Enabled {
		return fmt.Sprintf("%d", n)
	}
	return AlwaysByteSize(n, 2)
}

// AlwaysByteSize formats bytesize without checking `Enabled` first
func AlwaysByteSize(n uint64, precision int) string {
	if n < 1024 { // Never a decimal for byte-unit
		return printer.Sprintf("%d %s", n, byteSizes[0])
	}

	var nf float64 = float64(n)
	labelIdx := 0
	for nf >= 1024.0 && labelIdx < len(byteSizes)-1 {
		nf /= 1024.0
		labelIdx++
	}

	return printer.Sprintf("%.[2]*[1]f %[3]s", nf, precision, byteSizes[labelIdx])
}
