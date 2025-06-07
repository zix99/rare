package logger

import (
	"bytes"
	"log"
	"testing"

	"github.com/zix99/rare/pkg/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	Println("Howdy")
	DeferLogs()
	Println("Buffered")
	Print("Hi")
	Printf("Hi %v", 22)
	ImmediateLogs()
}

func TestLogFatal(t *testing.T) {
	exits := 0
	defer testutil.RestoreGlobals()
	testutil.SwitchGlobal(&OsExit, func(code int) {
		exits++
	})

	Fatal(1, "boom")
	Fatalln(1, "boom2")
	Fatalf(1, "Boom %v", "there")

	assert.Equal(t, 3, exits)
}

func TestLogCapture(t *testing.T) {
	logBuffer = new(bytes.Buffer)
	logger = log.New(logBuffer, "", 0)
	defer resetLogger()

	Println("Hello")
	Print("there")
	Printf("bob %v", 22)

	assert.Equal(t, "Hello\nthere\nbob 22\n", logBuffer.String())
}
