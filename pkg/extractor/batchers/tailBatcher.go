package batchers

import (
	"rare/pkg/extractor"
	"time"

	"github.com/hpcloud/tail"
)

func TailLineToChan(sourceName string, lines chan *tail.Line, batchSize int) <-chan extractor.InputBatch {
	output := make(chan extractor.InputBatch)
	go func() {
		batch := make([]extractor.BString, 0, batchSize)
		var batchStart uint64 = 1
	MAIN_LOOP:
		for {
			select {
			case line := <-lines:
				if line == nil || line.Err != nil {
					break MAIN_LOOP
				}
				batch = append(batch, extractor.BString(line.Text))
				if len(batch) >= batchSize {
					output <- extractor.InputBatch{
						Batch:      batch,
						Source:     sourceName,
						BatchStart: batchStart,
					}
					batchStart += uint64(len(batch))
					batch = make([]extractor.BString, 0, batchSize)
				}
			case <-time.After(500 * time.Millisecond):
				// Since we're tailing, if we haven't received any line in a bit, lets flush what we have
				if len(batch) > 0 {
					output <- extractor.InputBatch{
						Batch:      batch,
						Source:     sourceName,
						BatchStart: batchStart,
					}
					batchStart += uint64(len(batch))
					batch = make([]extractor.BString, 0, batchSize)
				}
			}
		}
		close(output)
	}()
	return output
}
