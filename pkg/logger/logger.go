package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

// ErrLog is the logger that is controlled by this log controller
var logger *log.Logger
var logBuffer *bytes.Buffer

const logPrefix = "[Log] "

func init() {
	resetLogger()
}

// DeferLogs enables the log-buffer and defers any logs from printing to the screen
func DeferLogs() {
	if logBuffer == nil {
		logBuffer = new(bytes.Buffer)
		logger = log.New(logBuffer, logPrefix, 0)
	}
}

// ImmediateLogs flushes logs and puts logging back into immediate mode
func ImmediateLogs() {
	if logBuffer != nil {
		os.Stderr.Write(logBuffer.Bytes())
		logBuffer = nil
		resetLogger()
	}
}

func resetLogger() {
	logger = log.New(os.Stderr, logPrefix, 0)
}

func Fatalln(s interface{}) {
	logger.Fatalln(s)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalf(s string, args ...interface{}) {
	Fatalln(fmt.Sprintf(s, args...))
}

func Println(s interface{}) {
	logger.Println(s)
}

func Print(v ...interface{}) {
	logger.Print(v...)
}

func Printf(s string, args ...interface{}) {
	Println(fmt.Sprintf(s, args...))
}
