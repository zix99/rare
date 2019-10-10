package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombiningChannels(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)

	combined := CombineChannels(c1, c2)
	c1 <- "a"
	c2 <- "b"
	assert.Equal(t, "a", <-combined)
	assert.Equal(t, "b", <-combined)

	close(c1)
	close(c2)

	_, more := <-combined
	assert.False(t, more)
}
