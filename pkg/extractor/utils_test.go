package extractor

import (
	"io"

	"github.com/zix99/rare/pkg/readahead"
)

func unbatchMatches[T any](c <-chan []T) []T {
	matches := make([]T, 0)
	for batch := range c {
		matches = append(matches, batch...)
	}
	return matches
}

func matchSetToString(matches []Match) []string {
	ret := make([]string, len(matches))
	for i, match := range matches {
		ret[i] = match.Extracted
	}
	return ret
}

func convertReaderToBatches(sourceName string, reader io.Reader, batchSize int) <-chan InputBatch {
	out := make(chan InputBatch)
	ra := readahead.NewImmediate(reader, 128*1024)

	go func() {
		batch := make([]BString, 0, batchSize)
		var batchStart uint64 = 1

		for ra.Scan() {
			batch = append(batch, ra.Bytes())
			if len(batch) >= batchSize {
				out <- InputBatch{batch, sourceName, batchStart}
				batchStart += uint64(len(batch))
				batch = make([]BString, 0, batchSize)
			}
		}
		if len(batch) > 0 {
			out <- InputBatch{batch, sourceName, batchStart}
		}

		close(out)
	}()

	return out
}
