//go:build !(linux && cgo && pcre2)

package fastregex

import "regexp"

/*
The fallback exposes the re2/regexp go implementaiton in the
cases where we can't compile with PCRE support
*/

const Version = "re2"

type compiledRegexp struct {
	*regexp.Regexp
	groupNames map[string]int
}

var (
	_ CompiledRegexp = &compiledRegexp{}
	_ Regexp         = &compiledRegexp{}
)

func (s *compiledRegexp) CreateInstance() Regexp {
	return s
}

func (s *compiledRegexp) SubexpNameTable() map[string]int {
	return s.groupNames
}

func CompileEx(expr string, posix bool) (CompiledRegexp, error) {
	re, err := buildRegexp(expr, posix)
	if err != nil {
		return nil, err
	}
	return &compiledRegexp{
		re,
		createGroupNameTable(re),
	}, nil
}

func buildRegexp(expr string, posix bool) (*regexp.Regexp, error) {
	if posix {
		return regexp.CompilePOSIX(expr)
	}
	return regexp.Compile(expr)
}

func createGroupNameTable(re *regexp.Regexp) (ret map[string]int) {
	ret = make(map[string]int)
	for idx, name := range re.SubexpNames() {
		if name != "" {
			ret[name] = idx
		}
	}
	return
}
