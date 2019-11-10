package extractor

import (
	"io"
	"rare/pkg/readahead"
	"sync"
)

// CombineChannels combines multiple string channels into a single (unordered)
//  string channel
func CombineChannels(channels ...<-chan []BString) <-chan []BString {
	if channels == nil {
		return nil
	}
	if len(channels) == 1 {
		return channels[0]
	}

	out := make(chan []BString)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(subchan <-chan []BString) {
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
func ConvertReaderToStringChan(reader io.ReadCloser, batchSize int) <-chan []BString {
	out := make(chan []BString)
	ra := readahead.New(reader, 128*1024)

	go func() {
		defer reader.Close()
		SyncReadAheadToBatchChannel(ra, batchSize, out)
		close(out)
	}()

	return out
}

func SyncReadAheadToBatchChannel(readahead *readahead.ReadAhead, batchSize int, out chan<- []BString) {
	batch := make([]BString, 0, batchSize)
	for {
		b := readahead.ReadLine()
		if b == nil {
			break
		}

		batch = append(batch, b)
		if len(batch) >= batchSize {
			out <- batch
			batch = make([]BString, 0, batchSize)
		}
	}
	if len(batch) > 0 {
		out <- batch
	}
}
