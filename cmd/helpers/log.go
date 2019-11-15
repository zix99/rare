package helpers

import (
	"bytes"
	"log"
	"os"
)

var ErrLog *log.Logger

const logPrefix = "[Log] "

var logBuffer *bytes.Buffer

func resetLog() {
	ErrLog = log.New(os.Stderr, "[Log] ", 0)
}

func init() {
	resetLog()
}

func EnableLogBuffer() {
	if logBuffer == nil {
		logBuffer = new(bytes.Buffer)
		ErrLog = log.New(logBuffer, logPrefix, 0)
	}
}

func DisableAndFlushLogBuffer() {
	if logBuffer != nil {
		os.Stderr.Write(logBuffer.Bytes())
		logBuffer = nil
		resetLog()
	}
}
