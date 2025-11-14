package batchers

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zix99/rare/pkg/extractor"
	"github.com/zix99/rare/pkg/humanize"
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/readahead"
)

// AutoFlushTimeout sets time before an auto-flushing reader will write a batch
const AutoFlushTimeout = 250 * time.Millisecond

type Batcher struct {
	c chan extractor.InputBatch

	// All mutex protected fields
	mux         sync.Mutex
	sourceCount int
	readCount   int
	errorCount  int
	activeFiles []string

	startTime, stopTime time.Time

	// Atomic fields (only used to compute performance metrics)
	readBytes uint64

	// Used only in StatusString to compute read rate
	lastRateUpdate          time.Time
	lastRate, lastRateBytes uint64
}

func newBatcher(bufferSize int) *Batcher {
	return &Batcher{
		c:              make(chan extractor.InputBatch, bufferSize),
		lastRateUpdate: time.Now(),
		startTime:      time.Now(),
	}
}

func (s *Batcher) close() {
	close(s.c)

	s.mux.Lock()
	s.stopTime = time.Now()
	s.mux.Unlock()
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

func (s *Batcher) ReadFiles() int {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.readCount
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
// [9/10 !1] 1.41 GB in 2.7s (~526.71 MB/s) | a b c (and 2 more...)
func (s *Batcher) StatusString() string {
	var sb strings.Builder
	sb.Grow(100)
	const maxFilesToWrite = 2

	s.mux.Lock()
	defer s.mux.Unlock()

	// Total files read
	if s.sourceCount > 1 {
		sb.WriteString(fmt.Sprintf("[%d/%d", s.readCount, s.sourceCount))
		// Errors
		if s.errorCount > 0 {
			sb.WriteString(fmt.Sprintf(" !%d", s.errorCount))
		}
		sb.WriteString("] ")
	}

	// Total read bytes
	readBytes := atomic.LoadUint64(&s.readBytes)
	sb.WriteString(humanize.ByteSize(readBytes))

	// Elapsed time
	elapsed := s.elapsedTimeNoLock()
	sb.WriteString(" in " + durationToString(elapsed))

	// Read rate
	if s.stopTime.IsZero() {
		// Progress
		elapsedTime := time.Since(s.lastRateUpdate).Seconds()
		if elapsedTime >= 0.5 {
			s.lastRate = uint64(float64(s.readBytes-s.lastRateBytes) / elapsedTime)
			s.lastRateBytes = s.readBytes
			s.lastRateUpdate = time.Now()
		}

		sb.WriteString(" (" + humanize.ByteSize(s.lastRate) + "/s)")
	} else {
		// Final
		rate := uint64(float64(readBytes) / elapsed.Seconds())
		sb.WriteString(" (~" + humanize.ByteSize(rate) + "/s)")
	}

	// Current actively read files
	writeFiles := min(len(s.activeFiles), maxFilesToWrite)
	if writeFiles > 0 {
		sb.WriteString(" | ")
		sb.WriteString(strings.Join(s.activeFiles[:writeFiles], ", "))

		if len(s.activeFiles) > maxFilesToWrite {
			sb.WriteString(fmt.Sprintf(" (and %d more...)", len(s.activeFiles)-maxFilesToWrite))
		}
	}

	return sb.String()
}

func (s *Batcher) elapsedTimeNoLock() time.Duration {
	if s.stopTime.IsZero() {
		return time.Since(s.startTime)
	}
	return s.stopTime.Sub(s.startTime)
}

// Variable duration pretty-printing
// Optimize to prevent terminal stutter/length changes (eg 2.1 2.11...)
func durationToString(d time.Duration) string {
	switch {
	case d < time.Second:
		return fmt.Sprintf("%03dms", d.Milliseconds())
	case d < time.Minute:
		return fmt.Sprintf("%.02fs", d.Truncate(10*time.Millisecond).Seconds())
	case d < time.Hour:
		return fmt.Sprintf("%dm%.1fs", int(d.Truncate(time.Minute).Minutes()), (d % time.Minute).Seconds())
	default:
		return d.Truncate(time.Second).String()
	}
}

// syncReaderToBatcher reads a reader buffer and breaks up its scans to `batchSize`
//
//	and writes the batch-sized results to a channel
func (s *Batcher) syncReaderToBatcher(sourceName string, reader io.Reader, batchSize, bufSize int) {
	readerMetrics := newReaderMetrics(reader)
	readahead := readahead.NewImmediate(readerMetrics, bufSize)
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
// it will flush the batch if n time has elapsed since the last flush, regardless of how many items are in the current batch
// Good for potentially slow or more interactive workloads (tail, stdin, etc)
func (s *Batcher) syncReaderToBatcherWithTimeFlush(sourceName string, reader io.Reader, batchSize, bufSize int, autoFlush time.Duration) {
	readerMetrics := newReaderMetrics(reader)
	readahead := readahead.NewImmediate(readerMetrics, bufSize)
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
