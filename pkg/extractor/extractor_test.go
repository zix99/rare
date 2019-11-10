package extractor

import (
	"regexp"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type StringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func TestSliceAssumptions(t *testing.T) {
	b := []byte("hello")
	z := b[1:]
	sr := (*string)(unsafe.Pointer(&z))
	b[1] = 'a'
	assert.Equal(t, "allo", *sr)
}

func BenchmarkRegexWithString(b *testing.B) {
	r := regexp.MustCompile("a(.*)")
	for n := 0; n < b.N; n++ {
		r.FindStringSubmatchIndex("abcdefg")
	}
}

func BenchmarkRegexWithBytes(b *testing.B) {
	r := regexp.MustCompile("a(.*)")
	val := []byte("abcdefg")
	for n := 0; n < b.N; n++ {
		r.FindSubmatchIndex(val)
	}
}
