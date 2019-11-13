package extractor

import (
	"io/ioutil"
	"regexp"
	"strings"
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

var testData = `abc 123
def 245
qqq 123
xxx`

func TestBasicExtractor(t *testing.T) {
	input := ConvertReaderToStringChan(ioutil.NopCloser(strings.NewReader(testData)), 1)
	ex, err := New(input, &Config{
		Regex:   `(\d+)`,
		Extract: "val:{1}",
		Workers: 1,
	})
	assert.NoError(t, err)

	val := <-ex.ReadChan()
	assert.Equal(t, "abc 123", val[0].Line)
	assert.Equal(t, 2, len(val[0].Groups))
	assert.Equal(t, 4, len(val[0].Indices))
	assert.Equal(t, "123", val[0].Groups[0])
	assert.Equal(t, "val:123", val[0].Extracted)
	assert.Equal(t, uint64(1), val[0].LineNumber)
	assert.Equal(t, uint64(1), val[0].MatchNumber)

	for _ = range val {
	} // Loop until closed

	assert.Equal(t, uint64(0), ex.IgnoredLines())
	assert.Equal(t, uint64(3), ex.MatchedLines())
	assert.Equal(t, uint64(4), ex.ReadLines())
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
