package readahead

import (
	"bufio"
	"rare/pkg/testutil"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestBasicReadingShortBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewBuffered(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingLongBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewBuffered(r, 1024)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingMidBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := NewBuffered(r, 20) // Just enough to read first line, but not both
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingNoNewTerm(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2")
	ra := NewBuffered(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestReadEmptyString(t *testing.T) {
	r := strings.NewReader("")
	ra := NewBuffered(r, 3)
	assert.Nil(t, ra.ReadLine())
}

func TestReadSingleCharString(t *testing.T) {
	r := strings.NewReader("A")
	ra := NewBuffered(r, 3)
	assert.Equal(t, []byte("A"), ra.ReadLine())
}

func TestBufferedDropCR(t *testing.T) {
	r := strings.NewReader("test\r\nthing")
	ra := NewBuffered(r, 3)
	assert.Equal(t, []byte("test"), ra.ReadLine())
	assert.Equal(t, []byte("thing"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestErrorHandling(t *testing.T) {
	errReader := iotest.TimeoutReader(strings.NewReader("Hello there you\nthis is a line\n"))
	ra := NewBuffered(errReader, 20)

	var hadError bool
	ra.OnError(func(e error) {
		hadError = true
	})

	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this"), ra.ReadLine()) // up to twentyith char
	assert.Nil(t, ra.ReadLine())
	assert.True(t, hadError)
}

func BenchmarkBuffered(b *testing.B) {
	r := testutil.NewTextGenerator(1024)
	ra := NewBuffered(r, 128*128)

	for i := 0; i < b.N; i++ {
		ra.Scan()
	}
}

func BenchmarkScanner(b *testing.B) {
	r := testutil.NewTextGenerator(1024)
	s := bufio.NewScanner(r)

	for i := 0; i < b.N; i++ {
		s.Scan()

		// Copy into a new memory slot as practically that's needed for how its consumed
		r := s.Bytes()
		data := make([]byte, len(r))
		copy(data, r)
	}
}
