// +build linux,cgo,pcre2

package fastregex

// PCRE2 Docs: https://www.pcre.org/current/doc/html/index.html

/*
#cgo LDFLAGS: -lpcre2-8
#cgo CFLAGS: -I/usr/include
#define PCRE2_CODE_UNIT_WIDTH 8
#include <pcre2.h>
#include <string.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

const Version = "libpcre2"

var JITEnabled = true

// readonly version of a compiled regexp
type pcre2Compiled struct {
	p          *C.pcre2_code
	groupCount int
	jitted     bool
}

// instance version
type pcre2Regexp struct {
	re *pcre2Compiled

	matchData *C.pcre2_match_data
	context   *C.pcre2_match_context
	jitStack  *C.pcre2_jit_stack

	ovec   *C.ulong // pointer to ovector within matchData
	goOvec []int    // Converted ovec
}

var _ Regexp = &pcre2Regexp{}

func Compile(expr string, posix bool) (CompiledRegexp, error) {
	if posix {
		return nil, errors.New("libpcre doesn't support posix")
	}

	bPtr := *(**C.uchar)(unsafe.Pointer(&expr))

	var errNum C.int
	var errOffset C.ulong
	compiled := C.pcre2_compile(bPtr, C.ulong(len([]byte(expr))), 0, &errNum, &errOffset, nil)
	if compiled == nil {
		buf := make([]C.uchar, 256)
		msgLen := C.pcre2_get_error_message(errNum, &buf[0], (C.size_t)(len(buf)))
		return nil, &compileError{
			Expr:    expr,
			Message: C.GoStringN((*C.char)(unsafe.Pointer(&buf[0])), msgLen),
			Offset:  int(errOffset),
		}
	}

	var jitted bool
	if JITEnabled {
		jitrc := C.pcre2_jit_compile(compiled, C.PCRE2_JIT_COMPLETE)
		if jitrc < 0 {
			C.pcre2_code_free(compiled)
			return nil, errors.New("pcre jit failed")
		}
		jitted = true
	}

	var groupCount C.uint32_t
	C.pcre2_pattern_info(compiled, C.PCRE2_INFO_CAPTURECOUNT, unsafe.Pointer(&groupCount))

	pcre := &pcre2Compiled{
		p:          compiled,
		groupCount: int(groupCount) + 1,
		jitted:     jitted,
	}
	runtime.SetFinalizer(pcre, func(f *pcre2Compiled) {
		C.pcre2_code_free(f.p)
	})
	return pcre, nil
}

func MustCompile(expr string, posix bool) CompiledRegexp {
	re, err := Compile(expr, posix)
	if err != nil {
		panic(err)
	}
	return re
}

func (s *pcre2Compiled) CreateInstance() Regexp {
	pcre := &pcre2Regexp{
		re: s,
	}

	if s.jitted {
		pcre.context = C.pcre2_match_context_create(nil)
		pcre.jitStack = C.pcre2_jit_stack_create(32*1024, 512*1024, nil)
		if pcre.jitStack == nil || pcre.context == nil {
			panic("pcre2: jit stack failure")
		}
		C.pcre2_jit_stack_assign(pcre.context, nil, unsafe.Pointer(pcre.jitStack))
	}

	pcre.matchData = C.pcre2_match_data_create(C.uint(s.groupCount), nil)
	if pcre.matchData == nil {
		panic("pcre2: match data failure")
	}
	pcre.ovec = C.pcre2_get_ovector_pointer(pcre.matchData)
	pcre.goOvec = make([]int, s.groupCount*2)

	runtime.SetFinalizer(pcre, func(f *pcre2Regexp) {
		if f.matchData != nil {
			C.pcre2_match_data_free(f.matchData)
		}
		if f.context != nil {
			C.pcre2_match_context_free(f.context)
		}
		if f.jitStack != nil {
			C.pcre2_jit_stack_free(f.jitStack)
		}
	})

	return pcre
}

func (s *pcre2Regexp) GroupCount() int {
	return s.re.groupCount
}

func (s *pcre2Regexp) Match(b []byte) bool {
	bPtr := (*C.uchar)(unsafe.Pointer(&b[0]))
	rc := C.pcre2_match(s.re.p, bPtr, C.size_t(len(b)), 0, 0, s.matchData, nil)
	return rc >= 0
}

func (s *pcre2Regexp) MatchString(str string) bool {
	bPtr := *(**C.uchar)(unsafe.Pointer(&str))
	rc := C.pcre2_match(s.re.p, bPtr, C.size_t(len(str)), 0, 0, s.matchData, s.context)
	return rc >= 0
}

// FindSubmatchIndex, like regexp, returns a set of string indecies where the results are
// FindSubmatchIndex is NOT thread-safe.  You need to create an instance of the fastregex engine
//  the return result is mutable-per response.  If you need after making a 2nd call to this function
//  a copy needs to be made
func (s *pcre2Regexp) FindSubmatchIndex(b []byte) []int {
	bPtr := (*C.uchar)(unsafe.Pointer(&b[0]))

	rc := C.pcre2_match(s.re.p, bPtr, C.size_t(len(b)), 0, 0, s.matchData, nil)
	if rc < 0 {
		return nil
	}

	for i := 0; i < s.re.groupCount*2; i++ {
		s.goOvec[i] = int(*(*C.ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(s.ovec)) + unsafe.Sizeof(*s.ovec)*uintptr(i))))
	}

	return s.goOvec
}

type compileError struct {
	Expr    string
	Message string
	Offset  int
}

var _ error = &compileError{}

func (s *compileError) Error() string {
	return fmt.Sprintf("Error in '%s', offset %d: %s", s.Expr, s.Offset, s.Message)
}
