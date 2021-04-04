// +build !linux !cgo !pcre1,!pcre2

package fastregex

import "regexp"

/*
The fallback exposes the re2/regexp go implementaiton in the
cases where we can't compile with PCRE support
*/

const Version = "re2"

type compiledRegexp struct {
	re *regexp.Regexp
}

var _ CompiledRegexp = &compiledRegexp{}

func (s *compiledRegexp) CreateInstance() Regexp {
	return s.re
}

func Compile(expr string) (CompiledRegexp, error) {
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}
	return &compiledRegexp{re}, nil
}

func MustCompile(expr string) CompiledRegexp {
	return &compiledRegexp{regexp.MustCompile(expr)}
}
