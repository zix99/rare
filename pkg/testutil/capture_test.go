package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdoutCapture(t *testing.T) {
	cap := NewCapture(&os.Stdout, false)
	fmt.Print("Hi thar!!!")
	cap.Close()

	assert.Equal(t, "Hi thar!!!", cap.String())
}

func TestStdinWrite(t *testing.T) {
	cap := NewCapture(&os.Stdin, true)
	cap.Writer().WriteString("inputstring\n")

	var s string
	fmt.Scanf("%s", &s)
	assert.Equal(t, "inputstring", s)

	cap.Close()
}

func TestCaptureFunc(t *testing.T) {
	stdout, stderr, err := Capture(func(w *os.File) error {
		fmt.Print("This is some output")
		fmt.Fprint(os.Stderr, "This is a log")
		w.WriteString("Hi")
		return nil
	})

	assert.Equal(t, "This is some output", stdout)
	assert.Equal(t, "This is a log", stderr)
	assert.NoError(t, err)
}
