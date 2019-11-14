package helpers

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"rare/pkg/extractor"
	"rare/pkg/readahead"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

const ReadAheadBufferSize = 128 * 1024

func tailLineToChan(lines chan *tail.Line, batchSize int) <-chan []extractor.BString {
	output := make(chan []extractor.BString)
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
					output <- batch
					batch = make([]extractor.BString, 0, batchSize)
				}
			case <-time.After(1000 * time.Millisecond):
				// Since we're tailing, if we haven't received any line in a bit, lets flush what we have
				if len(batch) > 0 {
					output <- batch
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
			stderrLog.Printf("Gunzip error for file %s: %v\n", filename, err)
		} else {
			file = zfile
		}
	}

	return file, nil
}

func openFilesToChan(filenames []string, gunzip bool, concurrency int, batchSize int) <-chan []extractor.BString {
	out := make(chan []extractor.BString, 128)
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
					stderrLog.Printf("Error opening file %s: %v\n", goFilename, err)
					return
				}
				defer file.Close()
				StartFileReading(goFilename)

				ra := readahead.New(file, ReadAheadBufferSize)
				extractor.SyncReadAheadToBatchChannel(ra, batchSize, out)

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

func isDir(path string) bool {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func globExpand(paths []string, recursive bool) []string {
	out := make([]string, 0)
	for _, p := range paths {
		if recursive && isDir(p) {
			filepath.Walk(p, func(walkPath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					out = append(out, walkPath)
				}
				return nil
			})
		} else {
			expanded, err := filepath.Glob(p)
			if err != nil {
				stderrLog.Printf("Path error: %v\n", err)
			} else {
				out = append(out, expanded...)
			}
		}
	}
	return out
}
