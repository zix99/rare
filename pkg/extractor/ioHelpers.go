package extractor

import (
	"bufio"
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
	// TODO: Use new readahead
	out := make(chan []BString)
	scanner := bufio.NewScanner(reader)
	bigBuf := make([]byte, 512*1024)
	scanner.Buffer(bigBuf, len(bigBuf))

	go func() {
		defer reader.Close()
		SyncScannerToBatchChannel(scanner, batchSize, out)
		close(out)
	}()

	return out
}

// SyncScannerToBatchChannel reads a scanner into []string chunks and writes to an output channel
func SyncScannerToBatchChannel(scanner *bufio.Scanner, batchSize int, out chan<- []BString) {
	batch := make([]BString, 0, batchSize)
	for scanner.Scan() {
		b := scanner.Bytes()
		cb := make(BString, len(b))
		copy(cb, b)

		batch = append(batch, cb)
		if len(batch) >= batchSize {
			out <- batch
			batch = make([]BString, 0, batchSize)
		}
	}
	if len(batch) > 0 {
		out <- batch
	}
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
