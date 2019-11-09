package helpers

import (
	"fmt"
	"strings"
	"sync"
)

var activeRead = []string{}
var activeReadMutex sync.Mutex
var sourceCount = 0
var readCount = 0

func IncSourceCount(delta int) {
	sourceCount += delta
}

func StartFileReading(source string) {
	activeReadMutex.Lock()
	activeRead = append(activeRead, source)
	activeReadMutex.Unlock()
}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
