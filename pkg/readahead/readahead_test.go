package readahead

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicReadingShortBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := New(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingLongBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := New(r, 1024)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingMidBuf(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2\n")
	ra := New(r, 20) // Just enough to read first line, but not both
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestBasicReadingNoNewTerm(t *testing.T) {
	r := strings.NewReader("Hello there you\nthis is line 2")
	ra := New(r, 3)
	assert.Equal(t, []byte("Hello there you"), ra.ReadLine())
	assert.Equal(t, []byte("this is line 2"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}

func TestReadEmptyString(t *testing.T) {
	r := strings.NewReader("")
	ra := New(r, 3)
	assert.Nil(t, ra.ReadLine())
}

func TestReadSingleCharString(t *testing.T) {
	r := strings.NewReader("A")
	ra := New(r, 3)
	assert.Equal(t, []byte("A"), ra.ReadLine())
}

func TestDropCR(t *testing.T) {
	r := strings.NewReader("test\r\nthing")
	ra := New(r, 3)
	assert.Equal(t, []byte("test"), ra.ReadLine())
	assert.Equal(t, []byte("thing"), ra.ReadLine())
	assert.Nil(t, ra.ReadLine())
}
