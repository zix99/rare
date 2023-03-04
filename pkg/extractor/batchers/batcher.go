package batchers

import (
	"fmt"
	"io"
	"rare/pkg/extractor"
	"rare/pkg/humanize"
	"rare/pkg/logger"
	"rare/pkg/readahead"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ReadAheadBufferSize is the default size of the read-ahead buffer
const ReadAheadBufferSize = 128 * 1024

// AutoFlushTimeout sets time before an auto-flushing reader will write a batch
const AutoFlushTimeout = 250 * time.Millisecond

type Batcher struct {
	c chan extractor.InputBatch

	mux         sync.Mutex
	sourceCount int
	readCount   int
	errorCount  int
	activeFiles []string

	readBytes               uint64
	lastRateUpdate          time.Time
	lastRate, lastRateBytes uint64
}

func newBatcher(bufferSize int) *Batcher {
	return &Batcher{
		c:              make(chan extractor.InputBatch, bufferSize),
		lastRateUpdate: time.Now(),
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

func (s *Batcher) incReadBytes(n uint64) {
	atomic.AddUint64(&s.readBytes, n)
}

func (s *Batcher) ReadBytes() uint64 {
	return atomic.LoadUint64(&s.readBytes)
}

func (s *Batcher) ReadErrors() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.errorCount
}

func (s *Batcher) ActiveFileCount() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return len(s.activeFiles)
}

// StatusString gets a formatted version of the current reader-set
func (s *Batcher) StatusString() string {
	var sb strings.Builder
	sb.Grow(100)
	const maxFilesToWrite = 2

	s.mux.Lock()
	// Total files read
	if s.sourceCount > 1 {
		sb.WriteString(fmt.Sprintf("[%d/%d] ", s.readCount, s.sourceCount))
	}

	// Rate / bytes
	readBytes := atomic.LoadUint64(&s.readBytes)
	sb.WriteString(humanize.ByteSize(readBytes) + " ")

	elapsedTime := time.Since(s.lastRateUpdate).Seconds()
	if elapsedTime >= 0.5 {
		s.lastRate = uint64(float64(s.readBytes-s.lastRateBytes) / elapsedTime)
		s.lastRateBytes = s.readBytes
		s.lastRateUpdate = time.Now()
	}

	sb.WriteString("(" + humanize.ByteSize(s.lastRate) + "/s) ")

	// Current actively read files
	writeFiles := min(len(s.activeFiles), maxFilesToWrite)
	if writeFiles > 0 {
		sb.WriteString("| ")
		sb.WriteString(strings.Join(s.activeFiles[:writeFiles], ", "))

		if len(s.activeFiles) > maxFilesToWrite {
			sb.WriteString(fmt.Sprintf(" (and %d more...)", len(s.activeFiles)-maxFilesToWrite))
		}
	}

	s.mux.Unlock()

	return sb.String()
}

// syncReaderToBatcher reads a reader buffer and breaks up its scans to `batchSize`
//  and writes the batch-sized results to a channel
func (s *Batcher) syncReaderToBatcher(sourceName string, reader io.Reader, batchSize int) {
	readerMetrics := newReaderMetrics(reader)
	readahead := readahead.NewImmediate(readerMetrics, ReadAheadBufferSize)
	readahead.OnError(func(e error) {
		s.incErrors()
		logger.Printf("Error reading %s: %v", sourceName, e)
	})

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

			s.incReadBytes(readerMetrics.CountReset())
		}
	}
	if len(batch) > 0 {
		s.c <- extractor.InputBatch{
			Batch:      batch,
			Source:     sourceName,
			BatchStart: batchStart,
		}
		s.incReadBytes(readerMetrics.CountReset())
	}
}

// syncReaderToBatcherWithTimeFlush is similar to `syncReaderToBatcher`, except if it gets a new line
// it will flush the batch if n time has elapsed since the last flush, irregardless of how many items are in the current batch
// Good for potentially slow or more interactive workloads (tail, stdin, etc)
func (s *Batcher) syncReaderToBatcherWithTimeFlush(sourceName string, reader io.Reader, batchSize int, autoFlush time.Duration) {
	readerMetrics := newReaderMetrics(reader)
	readahead := readahead.NewImmediate(readerMetrics, ReadAheadBufferSize)
	readahead.OnError(func(e error) {
		s.incErrors()
		logger.Printf("Error reading %s: %v", sourceName, e)
	})

	batch := make([]extractor.BString, 0, batchSize)
	var batchStart uint64 = 1
	lastBatchFlush := time.Now()

	for readahead.Scan() {
		batch = append(batch, readahead.Bytes())
		if len(batch) >= batchSize || time.Since(lastBatchFlush) >= autoFlush {
			s.c <- extractor.InputBatch{
				Batch:      batch,
				Source:     sourceName,
				BatchStart: batchStart,
			}
			batchStart += uint64(len(batch))
			batch = make([]extractor.BString, 0, batchSize)

			s.incReadBytes(readerMetrics.CountReset())
			lastBatchFlush = time.Now()
		}
	}
	if len(batch) > 0 {
		s.c <- extractor.InputBatch{
			Batch:      batch,
			Source:     sourceName,
			BatchStart: batchStart,
		}
		s.incReadBytes(readerMetrics.CountReset())
	}
}
