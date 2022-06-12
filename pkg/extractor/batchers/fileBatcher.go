package batchers

import (
	"compress/gzip"
	"io"
	"os"
	"rare/pkg/logger"
	"sync"
)

// openFilesToChan takes an iterated channel of filenames, options, and loads them all with
//  a max concurrency.  Returns a channel that will populate with input batches
func OpenFilesToChan(filenames <-chan string, gunzip bool, concurrency int, batchSize int) *Batcher {
	out := newBatcher(128)
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
			out.setSourceCount(readCount + len(bufferedFilenames))

			go func(goFilename string) {
				defer func() {
					<-sema
					wg.Done()
					out.stopFileReading(goFilename)
				}()

				var file io.ReadCloser
				file, err := openFileToReader(goFilename, gunzip)
				if err != nil {
					logger.Printf("Error opening file %s: %v", goFilename, err)
					out.incErrors()
					return
				}
				defer file.Close()

				out.startFileReading(goFilename)
				out.syncReaderToBatcher(goFilename, file, batchSize)
			}(filename)
		}

		wg.Wait()
		out.close()
	}()

	return out
}

func openFileToReader(filename string, gunzip bool) (io.ReadCloser, error) {
	baseFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	var file io.ReadCloser = baseFile

	if gunzip {
		zfile, err := gzip.NewReader(file)
		if err != nil {
			logger.Printf("Gunzip error for file %s: %v; Reading as plain file", filename, err)
			baseFile.Seek(0, io.SeekStart) // Rewind, since it probably took a few bytes to figure out this wasn't a gzip file
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
