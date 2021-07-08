package batchers

import (
	"fmt"
	"io"
	"rare/pkg/extractor"
	"rare/pkg/logger"
	"rare/pkg/readahead"
	"strings"
	"sync"
)

// ReadAheadBufferSize is the default size of the read-ahead buffer
const ReadAheadBufferSize = 128 * 1024

type Batcher struct {
	c chan extractor.InputBatch

	mux         sync.Mutex
	sourceCount int
	readCount   int
	errorCount  int
	activeFiles []string
}

func newBatcher(bufferSize int) *Batcher {
	return &Batcher{
		c: make(chan extractor.InputBatch, bufferSize),
	}
}

func (s *Batcher) close() {
	close(s.c)
}

func (s *Batcher) BatchChan() <-chan extractor.InputBatch {
	return s.c
}

// SetSourceCount sets the number of source files
func (s *Batcher) setSourceCount(count int) {
	s.sourceCount = count
}

// StartFileReading registers a given source as being read in the global read-pool
func (s *Batcher) startFileReading(source string) {
	s.mux.Lock()
	s.activeFiles = append(s.activeFiles, source)
	s.mux.Unlock()
}

// StopFileReading recognizes a source has stopped reading, and increments the fully-read counter
func (s *Batcher) stopFileReading(source string) {
	s.mux.Lock()
	for idx, ele := range s.activeFiles {
		if ele == source {
			s.activeFiles = append(s.activeFiles[:idx], s.activeFiles[idx+1:]...)
			s.readCount++
			break
		}
	}
	s.mux.Unlock()
}

func (s *Batcher) incErrors() {
	s.mux.Lock()
	s.errorCount++
	s.mux.Unlock()
}

// GetReadFileString gets a formatted version of the current reader-set
func (s *Batcher) StatusString() string {
	var sb strings.Builder
	const maxFilesToWrite = 2

	s.mux.Lock()
	if s.sourceCount > 1 && s.readCount != s.sourceCount {
		sb.WriteString(fmt.Sprintf("[%d/%d] ", s.readCount, s.sourceCount))
	}

	writeFiles := min(len(s.activeFiles), maxFilesToWrite)
	sb.WriteString(strings.Join(s.activeFiles[:writeFiles], ", "))

	if len(s.activeFiles) > maxFilesToWrite {
		sb.WriteString(fmt.Sprintf(" (and %d more...)", len(s.activeFiles)-maxFilesToWrite))
	}
	s.mux.Unlock()

	return sb.String()
}

func (s *Batcher) ReadErrors() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.errorCount
}

// syncReaderToBatcher reads a reader buffer and breaks up its scans to `batchSize`
//  and writes the batch-sized results to a channel
func (s *Batcher) syncReaderToBatcher(sourceName string, reader io.Reader, batchSize int) {
	readahead := readahead.New(reader, ReadAheadBufferSize)
	readahead.OnError = func(e error) {
		s.incErrors()
		logger.Printf("Error reading %s: %v", sourceName, e)
	}

	batch := make([]extractor.BString, 0, batchSize)
	var batchStart uint64 = 1
	for readahead.Scan() {
		batch = append(batch, readahead.Bytes())
		if len(batch) >= batchSize {
			s.c <- extractor.InputBatch{
				Batch:      batch,
				Source:     sourceName,
				BatchStart: batchStart,
			}
			batchStart += uint64(len(batch))
			batch = make([]extractor.BString, 0, batchSize)
		}
	}
	if len(batch) > 0 {
		s.c <- extractor.InputBatch{
			Batch:      batch,
			Source:     sourceName,
			BatchStart: batchStart,
		}
	}
}
