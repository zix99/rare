package extractor

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombiningChannels(t *testing.T) {
	c1 := make(chan []BString)
	c2 := make(chan []BString)

	combined := CombineChannels(c1, c2)
	c1 <- []BString{BString("a")}
	c2 <- []BString{BString("b")}
	assert.Equal(t, []BString{BString("a")}, <-combined)
	assert.Equal(t, []BString{BString("b")}, <-combined)

	close(c1)
	close(c2)

	_, more := <-combined
	assert.False(t, more)
}

func TestConvertReaderToStringChan(t *testing.T) {
	buf := ioutil.NopCloser(strings.NewReader("line1\nline2\nline3"))
	c := ConvertReaderToStringChan(buf, 100)
	batch := <-c
	assert.Equal(t, 3, len(batch))
	assert.Equal(t, BString("line1"), batch[0])
	assert.Equal(t, BString("line2"), batch[1])
	assert.Equal(t, BString("line3"), batch[2])
}

func TestConvertReaderToStringChanSmallBatch(t *testing.T) {
	buf := ioutil.NopCloser(strings.NewReader("line1\nline2\nline3"))
	c := ConvertReaderToStringChan(buf, 1)
	assert.Equal(t, BString("line1"), (<-c)[0])
	assert.Equal(t, BString("line2"), (<-c)[0])
	assert.Equal(t, BString("line3"), (<-c)[0])

	_, more := <-c
	assert.False(t, more)
}
