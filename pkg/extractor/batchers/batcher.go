package batchers

import (
	"fmt"
	"rare/pkg/extractor"
	"strings"
	"sync"
)

type Batcher struct {
	c chan extractor.InputBatch

	mux         sync.Mutex
	sourceCount int
	readCount   int
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
