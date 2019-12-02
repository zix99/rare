package helpers

import (
	"sync"
	"testing"
	"time"

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

	assert.Eventually(t, func() bool {
		return len(bc) == 100
	}, 1*time.Second, 10*time.Millisecond)
}
