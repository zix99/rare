package helpers

// Aggregate one channel into another, with a buffer
func bufferChan(in <-chan string, size int) <-chan string {
	out := make(chan string, size)
	go func() {
		for item := range in {
			out <- item
		}
		close(out)
	}()
	return out
}
