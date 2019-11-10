package extractor

import (
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
