package logger

import (
	"bytes"
	"log"
	"os"
	"sync"
)

// ErrLog is the logger that is controlled by this log controller
var logger *log.Logger
var logBuffer *bytes.Buffer
var mux sync.RWMutex

const logPrefix = "[Log] "

// Allow overriding exit for unit tests
var OsExit = os.Exit

func init() {
	resetLogger()
}

// DeferLogs enables the log-buffer and defers any logs from printing to the screen
func DeferLogs() {
	mux.Lock()
	defer mux.Unlock()

	if logBuffer == nil {
		logBuffer = new(bytes.Buffer)
		logger = log.New(logBuffer, logPrefix, 0)
	}
}

// ImmediateLogs flushes logs and puts logging back into immediate mode
func ImmediateLogs() {
	mux.Lock()
	defer mux.Unlock()

	if logBuffer != nil {
		os.Stderr.Write(logBuffer.Bytes())
		resetLogger()
	}
}

func resetLogger() {
	logBuffer = nil
	logger = log.New(os.Stderr, logPrefix, 0)
}

func Fatalln(code int, s interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Println(s)
	OsExit(code)
}

func Fatal(code int, v ...interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Print(v...)
	OsExit(code)
}

func Fatalf(code int, s string, args ...interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Printf(s, args...)
	OsExit(code)
}

func Println(s interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Println(s)
}

func Print(v ...interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Print(v...)
}

func Printf(s string, args ...interface{}) {
	mux.RLock()
	defer mux.RUnlock()

	logger.Printf(s, args...)
}
