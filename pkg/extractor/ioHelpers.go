package extractor

import (
	"rare/pkg/readahead"
)

// SyncReadAheadToBatchChannel reads a readahead buffer and breaks up its scans to `batchSize`
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
