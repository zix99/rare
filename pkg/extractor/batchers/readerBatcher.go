package batchers

import (
	"io"
)

func OpenReaderToChan(sourceName string, reader io.ReadCloser, batchSize int) *Batcher {
	out := newBatcher(128)

	go func() {
		defer reader.Close()
		defer out.close()
		out.startFileReading(sourceName)
		out.syncReaderToBatcherWithTimeFlush(sourceName, reader, batchSize, AutoFlushTimeout)
	}()

	return out
}
