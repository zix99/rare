package batchers

import (
	"io"
	"time"
)

func OpenReaderToChan(sourceName string, reader io.ReadCloser, batchSize int) *Batcher {
	out := newBatcher(128)

	go func() {
		defer reader.Close()
		defer out.close()
		out.startFileReading(sourceName)
		out.syncReaderToBatcherWithTimeFlush(sourceName, reader, batchSize, 250*time.Millisecond)
	}()

	return out
}
