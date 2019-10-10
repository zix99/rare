package extractor

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

// CombineChannels combines multiple string channels into a single (unordered)
//  string channel
func CombineChannels(channels ...chan string) chan string {
	if channels == nil {
		return nil
	}
	if len(channels) == 1 {
		return channels[0]
	}

	out := make(chan string, 2)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(subchan chan string) {
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
func ConvertReaderToStringChan(reader io.Reader) chan string {
	out := make(chan string)

	bufReader := bufio.NewReader(reader)
	go func() {
		for {
			line, err := bufReader.ReadString('\n')
			if err == io.EOF {
				break
			}
			out <- strings.TrimSuffix(line, "\n")
		}
		close(out)
	}()

	return out
}
