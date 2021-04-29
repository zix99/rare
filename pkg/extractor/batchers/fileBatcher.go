package batchers

import (
	"compress/gzip"
	"io"
	"os"
	"rare/pkg/extractor"
	"rare/pkg/extractor/readState"
	"rare/pkg/logger"
	"rare/pkg/readahead"
	"sync"
)

// ReadAheadBufferSize is the default size of the read-ahead buffer
const ReadAheadBufferSize = 128 * 1024

// openFilesToChan takes an iterated channel of filenames, options, and loads them all with
//  a max concurrency.  Returns a channel that will populate with input batches
func OpenFilesToChan(filenames <-chan string, gunzip bool, concurrency int, batchSize int) <-chan extractor.InputBatch {
	out := make(chan extractor.InputBatch, 128)
	sema := make(chan struct{}, concurrency)

	// Load as many files as the sema allows
	go func() {
		var wg sync.WaitGroup
		readCount := 0

		bufferedFilenames := bufferChan(filenames, 1000)
		for filename := range bufferedFilenames {
			sema <- struct{}{}

			wg.Add(1)
			readCount++
			readState.SetSourceCount(readCount + len(bufferedFilenames))

			go func(goFilename string) {
				defer func() {
					<-sema
					wg.Done()
					readState.StopFileReading(goFilename)
				}()

				var file io.ReadCloser
				file, err := openFileToReader(goFilename, gunzip)
				if err != nil {
					logger.Printf("Error opening file %s: %v", goFilename, err)
					return
				}
				defer file.Close()
				readState.StartFileReading(goFilename)

				ra := readahead.New(file, ReadAheadBufferSize)
				ra.OnError = func(e error) {
					logger.Printf("Error reading %s: %v", goFilename, e)
				}
				extractor.SyncReadAheadToBatchChannel(goFilename, ra, batchSize, out)
			}(filename)
		}

		wg.Wait()
		close(out)
	}()

	return out
}

func openFileToReader(filename string, gunzip bool) (io.ReadCloser, error) {
	var file io.ReadCloser
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if gunzip {
		zfile, err := gzip.NewReader(file)
		if err != nil {
			logger.Printf("Gunzip error for file %s: %v", filename, err)
		} else {
			file = zfile
		}
	}

	return file, nil
}

// Aggregate one channel into another, with a buffer
func bufferChan(in <-chan string, size int) <-chan string {
	out := make(chan string, size)
	go func() {
		for item := range in {
			out <- item
		}
		close(out)
	}()
	return out
}
