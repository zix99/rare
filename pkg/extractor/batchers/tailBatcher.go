package batchers

import (
	"rare/pkg/extractor"
	"rare/pkg/logger"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

// TailFilesToChan tails a set of files to an input batcher that can be consumed by extractor
//  unlike a normal file batcher, this will attempt to tail all files at once
func TailFilesToChan(filenames <-chan string, batchSize int, reopen, poll bool) *Batcher {
	out := newBatcher(128)

	go func() {
		var wg sync.WaitGroup

		for filename := range filenames {
			wg.Add(1)
			go func(filename string) {
				defer func() {
					wg.Done()
					out.stopFileReading(filename)
				}()

				out.startFileReading(filename)
				fileTail, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: reopen, Poll: poll})
				if err != nil {
					logger.Print("Unable to open file: ", err)
					out.incErrors()
					return
				}

				err = tailLineToChan(filename, fileTail.Lines, batchSize, out.c)
				if err != nil {
					logger.Print("Error tailing file: ", err)
					out.incErrors()
				}
			}(filename)
		}

		wg.Wait()
		out.close()
	}()

	return out
}

func tailLineToChan(sourceName string, lines <-chan *tail.Line, batchSize int, output chan<- extractor.InputBatch) (err error) {
	batch := make([]extractor.BString, 0, batchSize)
	var batchStart uint64 = 1

MAIN_LOOP:
	for {
		select {
		case line := <-lines:
			if line == nil {
				break MAIN_LOOP
			}
			if line.Err != nil {
				err = line.Err
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
	return
}
