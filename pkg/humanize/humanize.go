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
