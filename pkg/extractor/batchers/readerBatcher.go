package batchers

import (
	"io"
	"rare/pkg/extractor"
	"rare/pkg/readahead"
)

func OpenReaderToChan(sourceName string, reader io.ReadCloser, batchSize int) *Batcher {
	out := newBatcher(128)
	ra := readahead.New(reader, 128*1024)

	go func() {
		defer reader.Close()
		defer out.close()
		out.startFileReading(sourceName)
		extractor.SyncReadAheadToBatchChannel(sourceName, ra, batchSize, out.c)
	}()

	return out
}
