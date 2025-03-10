package termformat

import (
	"rare/pkg/humanize"
	"strconv"
)

type Formatter func(val, min, max int64) string

type Simple func(val int64) string

// Passthru raw-formats the int (strconv)
func Passthru(val, min, max int64) string {
	return strconv.FormatInt(val, 10)
}

// Default formatter passes-through to humanize
func Default(val, min, max int64) string {
	return humanize.Hi(val)
}

// Maps a simple formatter (value-only) to a full formatter (value, min, max)
func ToFormatter(simple Simple) Formatter {
	return func(val, min, max int64) string {
		return simple(val)
	}
}
