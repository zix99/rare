package extractor

import (
	"regexp"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestSliceAssumptions(t *testing.T) {
	b := []byte("hello")
	z := b[1:]
	sr := *(*string)(unsafe.Pointer(&z))
	sr2 := sr[0:1]
	b[1] = 'a'
	assert.Equal(t, "allo", sr)
	assert.Equal(t, "a", sr2)
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
