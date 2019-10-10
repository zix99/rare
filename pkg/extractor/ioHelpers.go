package extractor

import (
	"bufio"
	"io"
	"strings"
	"sync"
)

func CombineChannels(channels ...chan string) chan string {
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
