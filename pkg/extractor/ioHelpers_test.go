package extractor

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombiningChannels(t *testing.T) {
	c1 := make(chan InputBatch)
	c2 := make(chan InputBatch)

	combined := CombineChannels(c1, c2)
	c1 <- InputBatch{[]BString{BString("a")}, "", 0}
	assert.Equal(t, []BString{BString("a")}, (<-combined).Batch)
	c2 <- InputBatch{[]BString{BString("b")}, "", 0}
	assert.Equal(t, []BString{BString("b")}, (<-combined).Batch)

	close(c1)
	close(c2)

	_, more := <-combined
	assert.False(t, more)
}

func TestConvertReaderToStringChan(t *testing.T) {
	buf := ioutil.NopCloser(strings.NewReader("line1\nline2\nline3"))
	c := ConvertReaderToStringChan("src", buf, 100)
	batch := <-c
	assert.Equal(t, "src", batch.Source)
	assert.Equal(t, 3, len(batch.Batch))
	assert.Equal(t, BString("line1"), batch.Batch[0])
	assert.Equal(t, BString("line2"), batch.Batch[1])
	assert.Equal(t, BString("line3"), batch.Batch[2])
}

func TestConvertReaderToStringChanSmallBatch(t *testing.T) {
	buf := ioutil.NopCloser(strings.NewReader("line1\nline2\nline3"))
	c := ConvertReaderToStringChan("src", buf, 1)
	assert.Equal(t, BString("line1"), (<-c).Batch[0])
	assert.Equal(t, BString("line2"), (<-c).Batch[0])
	assert.Equal(t, BString("line3"), (<-c).Batch[0])

	_, more := <-c
	assert.False(t, more)
}
