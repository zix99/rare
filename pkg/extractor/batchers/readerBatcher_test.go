package batchers

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenReaderToChan(t *testing.T) {
	r := ioutil.NopCloser(strings.NewReader("Hello\nthere\nbob"))
	b := OpenReaderToChan("src", r, 1)

	b1 := <-b.BatchChan()
	assert.Equal(t, "src", b1.Source)
	assert.Equal(t, "Hello", string(b1.Batch[0]))
	assert.Equal(t, uint64(1), b1.BatchStart)

	assert.Equal(t, 1, b.ActiveFileCount())
}
