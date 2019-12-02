package helpers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferingChan(t *testing.T) {
	var wg sync.WaitGroup

	c := make(chan string)
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			c <- "hi"
		}
		close(c)
		wg.Done()
	}()

	bc := bufferChan(c, 100)
	wg.Wait()
	assert.Equal(t, 100, len(bc))
}
