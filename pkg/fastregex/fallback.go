// +build !linux nopcre

package fastregex

import "regexp"

/*
The fallback exposes the re2/regexp go implementaiton in the
cases where we can't compile with PCRE support
*/

const Version = "re2"

func Compile(expr string) (Regexp, error) {
	return regexp.Compile(expr)
}

func MustCompile(expr string) Regexp {
	return regexp.MustCompile(expr)
}
