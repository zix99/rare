package extractor

import (
	"bufio"
	"io"
	"sync"
)

type semiLock struct{}

// CombineChannels combines multiple string channels into a single (unordered)
//  string channel
func CombineChannels(channels ...<-chan []string) <-chan []string {
	if channels == nil {
		return nil
	}
	if len(channels) == 1 {
		return channels[0]
	}

	out := make(chan []string)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(subchan <-chan []string) {
			for {
				s, more := <-subchan
				if !more {
					break
				}
				out <- s
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ConvertReaderToStringChan converts an io.reader to a string channel
//  where it's separated by a new-line
func ConvertReaderToStringChan(reader io.ReadCloser) <-chan []string {
	out := make(chan []string)
	scanner := bufio.NewScanner(reader)
	bigBuf := make([]byte, 512*1024)
	scanner.Buffer(bigBuf, len(bigBuf))

	go func() {
		defer reader.Close()
		for scanner.Scan() {
			out <- []string{scanner.Text()} // TODO: Batching
		}
		close(out)
	}()

	return out
}
