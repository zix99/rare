package readState

import (
	"fmt"
	"strings"
	"sync"
)

var activeRead = []string{}
var activeReadMutex sync.Mutex
var sourceCount = 0
var readCount = 0

// IncSourceCount increments the count of number of sources by delta
func IncSourceCount(delta int) {
	sourceCount += delta
}

// SetSourceCount sets the number of source files
func SetSourceCount(count int) {
	sourceCount = count
}

// StartFileReading registers a given source as being read in the global read-pool
func StartFileReading(source string) {
	activeReadMutex.Lock()
	activeRead = append(activeRead, source)
	activeReadMutex.Unlock()
}

// StopFileReading recognizes a source has stopped reading, and increments the fully-read counter
func StopFileReading(source string) {
	activeReadMutex.Lock()
	for idx, ele := range activeRead {
		if ele == source {
			activeRead = append(activeRead[:idx], activeRead[idx+1:]...)
			break
		}
	}
	readCount++
	activeReadMutex.Unlock()
}

// GetReadFileString gets a formatted version of the current reader-set
func GetReadFileString() string {
	var sb strings.Builder
	const maxFilesToWrite = 2

	activeReadMutex.Lock()
	if sourceCount > 1 && readCount != sourceCount {
		sb.WriteString(fmt.Sprintf("[%d/%d] ", readCount, sourceCount))
	}

	writeFiles := min(len(activeRead), maxFilesToWrite)
	sb.WriteString(strings.Join(activeRead[:writeFiles], ", "))

	if len(activeRead) > maxFilesToWrite {
		sb.WriteString(fmt.Sprintf(" (and %d more...)", len(activeRead)-maxFilesToWrite))
	}
	activeReadMutex.Unlock()

	return sb.String()
}
