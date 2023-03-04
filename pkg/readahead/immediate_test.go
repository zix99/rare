package readahead

import (
	"rare/pkg/testutil"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestImmediateBasicReadingShortBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewImmediate(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateBasicReadingLongBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewImmediate(r, 1024)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateBasicReadingMidBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewImmediate(r, 20) // Just enough to read first line, but not both
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateBasicReadingNoNewTerm(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2")
	ra := NewImmediate(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateReadEmptyString(t *testing.T) {
	r := strings.NewReader("")
	ra := NewImmediate(r, 3)
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateReadSingleCharString(t *testing.T) {
	r := strings.NewReader("A")
	ra := NewImmediate(r, 3)
	assert.Equal(t, []byte("A"), ra.ReadLine())
}

func TestImmediateDropCR(t *testing.T) {
	r := strings.NewReader("test\r\nthing")
	ra := NewImmediate(r, 3)
	assert.Equal(t, []byte("test"), ra.ReadLine())
	assert.Equal(t, []byte("thing"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestImmediateErrorHandling(t *testing.T) {
	errReader := iotest.TimeoutReader(strings.NewReader("Hello there you\nthis is a line\n"))
	ra := NewImmediate(errReader, 20)

	var hadError bool
	ra.OnError(func(e error) {
		hadError = true
	})

	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this"), ra.ReadLine()) // up to twentyith char
	assert.Nil(t, ra.ReadLine())
	assert.True(t, hadError)
}

func BenchmarkImmediate(b *testing.B) {
	r := testutil.NewTextGenerator(1024)
	ra := NewImmediate(r, 128*128)

	for i := 0; i < b.N; i++ {
		ra.Scan()
	}
}
