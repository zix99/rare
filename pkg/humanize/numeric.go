package humanize

import (
	"strconv"
	"strings"
)

const (
	baseSeparator    = ','
	decimalSeparator = '.'
)

type IntType interface {
	int64 | int32 | uint64 | uint32 | int | uint
}

func humanizeInt[T IntType](v T) string {
	var buf [32]byte // stack alloc

	if v >= 0 && v < 100 { // faster for small numbers
		return strconv.FormatInt(int64(v), 10)
	}

	negative := v < 0
	if negative {
		v = -v
	}

	ci := 0
	idx := len(buf) - 1
	for v > 0 {
		if ci == 3 {
			buf[idx] = baseSeparator
			ci = 0
			idx--
		}

		buf[idx] = byte('0' + (v % 10))
		idx--
		ci++
		v /= 10
	}

	if negative {
		buf[idx] = '-'
		idx--
	}

	return string(buf[idx+1:])
}

func humanizeFloat(v float64, decimals int) string {
	// Float to string is complicated, but can leverage FormatFload and insert commas
	s := strconv.FormatFloat(v, 'f', decimals, 64)

	if v > -1000.0 && v < 1000.0 {
		// performance escape hatch when no commas
		return s
	}

	dIdx := strings.IndexByte(s, '.')
	if dIdx < 0 { // no decimal
		dIdx = len(s)
	}
	negative := s[0] == '-'

	ret := make([]byte, 0, len(s)*2)

	// write base
	c3 := 3 - (dIdx % 3)

	for i := 0; i < dIdx; i++ {
		if c3 == 3 {
			if (!negative && i > 0) || (negative && i > 1) {
				ret = append(ret, baseSeparator)
			}
			c3 = 0
		}
		ret = append(ret, s[i])
		c3++
	}

	// write decimal
	if dIdx < len(s) {
		ret = append(ret, decimalSeparator)
		for i := dIdx + 1; i < dIdx+decimals+1 && i < len(s); i++ {
			ret = append(ret, s[i])
		}
	}

	return string(ret)
}
