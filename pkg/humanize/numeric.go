package humanize

import (
	"bytes"
	"math"
	"strconv"
)

/*
Previously, rare used the `message` i18n go library to add commas to numbers, but
as it turns out that was a bit overkill (Benchmarking shows easily 10x slower, and added 600 KB to the
binary).  In an effort to pull out and streamline simpler parts of the overall process,
the two below functions are implementations of the simplistic english-only localization
of numbers
*/

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

func isDigit(s byte) bool {
	return s >= '0' && s <= '9'
}

func humanizeFloat(v float64, decimals int) string {
	// Special cases
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 0) {
		return "Inf"
	}

	// Float to string is complicated, but can leverage FormatFload and insert commas
	var buf [64]byte // Operations on the stack
	s := strconv.AppendFloat(buf[:0], v, 'f', decimals, 64)

	if v > -1000.0 && v < 1000.0 {
		// performance escape hatch when no commas
		return string(s)
	}

	negative := s[0] == '-'
	if !isDigit(s[0]) { // assume it's a sign/prefix
		s = s[1:]
	}

	decIdx := bytes.IndexByte(s, '.')
	if decIdx < 0 { // no decimal
		decIdx = len(s)
	}

	// Return stack buf
	var retbuf [64]byte
	ret := retbuf[:0]

	if negative {
		ret = append(ret, '-')
	}

	// write base
	c3 := 3 - (decIdx % 3)

	for i := 0; i < decIdx; i++ {
		if c3 == 3 {
			if i > 0 {
				ret = append(ret, baseSeparator)
			}
			c3 = 0
		}
		ret = append(ret, s[i])
		c3++
	}

	// write decimal
	if decIdx < len(s) {
		ret = append(ret, decimalSeparator)
		ret = append(ret, s[decIdx+1:]...)
	}

	return string(ret)
}
