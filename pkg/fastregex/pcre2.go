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
	groupNames map[string]int
	jitted     bool
}

var _ CompiledRegexp = &pcre2Compiled{}

// instance version
type pcre2Regexp struct {
	re *pcre2Compiled

	matchData *C.pcre2_match_data
	context   *C.pcre2_match_context
	jitStack  *C.pcre2_jit_stack

	ovec *C.ulong // pointer to ovector within matchData
}

var _ Regexp = &pcre2Regexp{}

func CompileEx(expr string, posix bool) (CompiledRegexp, error) {
	if posix {
		return nil, errors.New("libpcre doesn't support posix")
	}

	bPtr := *(**C.uchar)(unsafe.Pointer(&expr))

	var errNum C.int
	var errOffset C.ulong
	compiled := C.pcre2_compile(bPtr, C.ulong(len(expr)), 0, &errNum, &errOffset, nil)
	if compiled == nil {
		buf := make([]C.uchar, 256)
		msgLen := C.pcre2_get_error_message(errNum, &buf[0], C.size_t(len(buf)))
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

	groupCount := pcre2PatternInfoUint32(compiled, C.PCRE2_INFO_CAPTURECOUNT)
	groupNames := pcre2BuildGroupNameTable(compiled)

	pcre := &pcre2Compiled{
		p:          compiled,
		groupCount: groupCount + 1,
		groupNames: groupNames,
		jitted:     jitted,
	}
	runtime.SetFinalizer(pcre, func(f *pcre2Compiled) {
		C.pcre2_code_free(f.p)
	})
	return pcre, nil
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

func (s *pcre2Regexp) SubexpNameTable() map[string]int {
	return s.re.groupNames
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
func (s *pcre2Regexp) FindSubmatchIndex(b []byte) []int {
	if len(b) == 0 {
		return nil
	}

	bPtr := (*C.uchar)(unsafe.Pointer(&b[0]))

	rc := C.pcre2_match(s.re.p, bPtr, C.size_t(len(b)), 0, 0, s.matchData, s.context)
	if rc < 0 {
		return nil
	}

	ret := make([]int, s.re.groupCount*2)
	for i := 0; i < s.re.groupCount*2; i++ {
		ret[i] = int(*(*C.ulong)(unsafe.Pointer(uintptr(unsafe.Pointer(s.ovec)) + unsafe.Sizeof(*s.ovec)*uintptr(i))))
	}

	return ret
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

func pcre2BuildGroupNameTable(p *C.pcre2_code) map[string]int {
	ret := make(map[string]int)

	nameCount := pcre2PatternInfoUint32(p, C.PCRE2_INFO_NAMECOUNT)
	if nameCount > 0 {
		nameEntrySize := pcre2PatternInfoUint32(p, C.PCRE2_INFO_NAMEENTRYSIZE)
		table := pcre2PatternInfoBytes(p, C.PCRE2_INFO_NAMETABLE, nameCount*nameEntrySize)

		for i := 0; i < nameCount; i++ {
			row := table[i*nameEntrySize : (i+1)*nameEntrySize]
			groupIndex := (int(row[0]) << 8) | int(row[1])
			name := nullTermString(row[2:])
			ret[name] = groupIndex
		}
	}

	return ret
}

func pcre2PatternInfoBytes(p *C.pcre2_code, code C.uint, size int) []byte {
	var ret *C.uchar
	C.pcre2_pattern_info(p, code, unsafe.Pointer(&ret))
	return C.GoBytes(unsafe.Pointer(ret), C.int(size))
}

func pcre2PatternInfoUint32(p *C.pcre2_code, code C.uint) int {
	var ret C.uint32_t
	C.pcre2_pattern_info(p, code, unsafe.Pointer(&ret))
	return int(ret)
}

func nullTermString(cstr []byte) string {
	for i := 0; i < len(cstr); i++ {
		if cstr[i] == 0 {
			return string(cstr[:i])
		}
	}
	return string(cstr)
}
