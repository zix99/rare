package batchers

import (
	"io"
)

func OpenReaderToChan(sourceName string, reader io.ReadCloser, batchSize, batchBuffer, readBufSize int) *Batcher {
	out := newBatcher(batchBuffer)

	go func() {
		defer reader.Close()
		defer out.close()
		out.startFileReading(sourceName)
		out.syncReaderToBatcherWithTimeFlush(sourceName, reader, batchSize, readBufSize, AutoFlushTimeout)
		out.stopFileReading(sourceName)
	}()

	return out
}
