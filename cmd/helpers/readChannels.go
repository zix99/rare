package helpers

import (
	"compress/gzip"
	"io"
	"os"
	"rare/pkg/extractor"
	"rare/pkg/readahead"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

// ReadAheadBufferSize is the default size of the read-ahead buffer
const ReadAheadBufferSize = 128 * 1024

func tailLineToChan(sourceName string, lines chan *tail.Line, batchSize int) <-chan extractor.InputBatch {
	output := make(chan extractor.InputBatch)
	go func() {
		batch := make([]extractor.BString, 0, batchSize)
	MAIN_LOOP:
		for {
			select {
			case line := <-lines:
				if line == nil || line.Err != nil {
					break MAIN_LOOP
				}
				batch = append(batch, extractor.BString(line.Text))
				if len(batch) >= batchSize {
					output <- extractor.InputBatch{batch, sourceName}
					batch = make([]extractor.BString, 0, batchSize)
				}
			case <-time.After(500 * time.Millisecond):
				// Since we're tailing, if we haven't received any line in a bit, lets flush what we have
				if len(batch) > 0 {
					output <- extractor.InputBatch{batch, sourceName}
					batch = make([]extractor.BString, 0, batchSize)
				}
			}
		}
		close(output)
	}()
	return output
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
			ErrLog.Printf("Gunzip error for file %s: %v\n", filename, err)
		} else {
			file = zfile
		}
	}

	return file, nil
}

func openFilesToChan(filenames []string, gunzip bool, concurrency int, batchSize int) <-chan extractor.InputBatch {
	out := make(chan extractor.InputBatch, 128)
	sema := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	wg.Add(len(filenames))
	IncSourceCount(len(filenames))

	// Load as many files as the sema allows
	go func() {
		for _, filename := range filenames {
			sema <- struct{}{}

			go func(goFilename string) {
				var file io.ReadCloser
				file, err := openFileToReader(goFilename, gunzip)
				if err != nil {
					ErrLog.Printf("Error opening file %s: %v\n", goFilename, err)
					return
				}
				defer file.Close()
				StartFileReading(goFilename)

				ra := readahead.New(file, ReadAheadBufferSize)
				ra.OnError = func(e error) {
					ErrLog.Printf("Error reading %s: %v\n", goFilename, e)
				}
				extractor.SyncReadAheadToBatchChannel(goFilename, ra, batchSize, out)

				<-sema
				wg.Done()
				StopFileReading(goFilename)
			}(filename)
		}
	}()

	// Wait on all files, and close chan
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
