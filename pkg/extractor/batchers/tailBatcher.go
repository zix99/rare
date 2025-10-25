package batchers

import (
	"sync"

	"github.com/zix99/rare/pkg/followreader"
	"github.com/zix99/rare/pkg/logger"
)

// TailFilesToChan tails a set of files to an input batcher that can be consumed by extractor
//
//	unlike a normal file batcher, this will attempt to tail all files at once
func TailFilesToChan(filenames <-chan string, batchSize, batchBuffer, readBufSize int, reopen, poll, tail bool) *Batcher {
	out := newBatcher(batchBuffer)

	go func() {
		var wg sync.WaitGroup
		for filename := range filenames {
			wg.Add(1)
			go func(filename string) {
				defer func() {
					wg.Done()
					out.stopFileReading(filename)
				}()

				r, err := followreader.New(filename, reopen, poll)
				if err != nil {
					logger.Print("Unable to open file: ", err)
					out.incErrors()
					return
				}
				if tail {
					if err := r.Drain(); err != nil {
						logger.Print("Unable to tail file source: ", err)
						out.incErrors()
					}
				}

				out.startFileReading(filename)

				out.syncReaderToBatcherWithTimeFlush(filename, r, batchSize, readBufSize, AutoFlushTimeout)
			}(filename)
		}

		wg.Wait()
		out.close()
	}()

	return out
}
