//go:build !(linux && cgo && pcre2)

package fastregex

import (
	"io"
	"regexp"
	_ "unsafe"
)

/*
The fallback exposes the re2/regexp go implementaiton in the
cases where we can't compile with PCRE support
*/

const Version = "re2"

type compiledRegexp struct {
	*regexp.Regexp
	groupNames map[string]int
	bufSize    int
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

//go:linkname regexp_doExecute regexp.(*Regexp).doExecute
func regexp_doExecute(*regexp.Regexp, io.RuneReader, []byte, string, int, int, []int) []int

func (s *compiledRegexp) FindSubmatchIndexDst(b []byte, dst []int) []int {
	// HACK: By accessing the underlying function of FindSubmatchIndex, we're able to avoid
	// an allocation done by the initial call, which seems to save 25-33% performance generally
	// and also later gc cleanups
	// Though hacky, this should be safe for a pinned version, and will have plenty of tests around it
	return regexp_doExecute(s.Regexp, nil, b, "", 0, s.bufSize, dst)
}

func (s *compiledRegexp) MatchBufSize() int {
	return s.bufSize
}

func CompileEx(expr string, posix bool) (CompiledRegexp, error) {
	re, err := buildRegexp(expr, posix)
	if err != nil {
		return nil, err
	}
	return &compiledRegexp{
		re,
		createGroupNameTable(re),
		(re.NumSubexp() + 1) * 2,
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
