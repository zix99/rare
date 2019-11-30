package extractor

import (
	"io"
	"rare/pkg/readahead"
	"sync"
)

// CombineChannels combines multiple string channels into a single (unordered)
//  string channel
func CombineChannels(channels ...<-chan InputBatch) <-chan InputBatch {
	if channels == nil {
		return nil
	}
	if len(channels) == 1 {
		return channels[0]
	}

	out := make(chan InputBatch)
	var wg sync.WaitGroup

	for _, c := range channels {
		wg.Add(1)
		go func(subchan <-chan InputBatch) {
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
func ConvertReaderToStringChan(sourceName string, reader io.ReadCloser, batchSize int) <-chan InputBatch {
	out := make(chan InputBatch)
	ra := readahead.New(reader, 128*1024)

	go func() {
		defer reader.Close()
		SyncReadAheadToBatchChannel(sourceName, ra, batchSize, out)
		close(out)
	}()

	return out
}

// SyncReadAheadToBatchChannel reads a readahead buffer and breaks up its scants to `batchSize`
//  and writes the batch-sized results to a channel
func SyncReadAheadToBatchChannel(sourceName string, readahead *readahead.ReadAhead, batchSize int, out chan<- InputBatch) {
	batch := make([]BString, 0, batchSize)
	var batchStart uint64 = 1
	for readahead.Scan() {
		batch = append(batch, readahead.Bytes())
		if len(batch) >= batchSize {
			out <- InputBatch{batch, sourceName, batchStart}
			batchStart += uint64(len(batch))
			batch = make([]BString, 0, batchSize)
		}
	}
	if len(batch) > 0 {
		out <- InputBatch{batch, sourceName, batchStart}
	}
}
