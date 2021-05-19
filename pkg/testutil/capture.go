package testutil

import (
	"bytes"
	"io"
	"os"
)

func Capture(f func(w *os.File) error) (stdout, stderr string, err error) {
	outcap := NewCapture(&os.Stdout, false)
	errcap := NewCapture(&os.Stderr, false)
	incap := NewCapture(&os.Stdin, true)

	retErr := f(incap.Writer())

	outcap.Close()
	errcap.Close()
	incap.Close()

	return outcap.String(), errcap.String(), retErr
}

type FileCapture struct {
	ptr            **os.File
	orig           *os.File
	reader, writer *os.File

	closeWait <-chan string
	result    string
}

func NewCapture(ptr **os.File, writer bool) *FileCapture {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	orig := *ptr
	var closeWait chan string

	if writer {
		*ptr = r
	} else {
		*ptr = w
		closeWait = make(chan string)
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			closeWait <- buf.String()
		}()
	}

	return &FileCapture{
		ptr:       ptr,
		orig:      orig,
		reader:    r,
		writer:    w,
		closeWait: closeWait,
	}
}

func (s *FileCapture) Close() {
	s.writer.Close()
	*s.ptr = s.orig
	if s.closeWait != nil {
		s.result = <-s.closeWait
	}
}

func (s *FileCapture) String() string {
	return s.result
}

func (s *FileCapture) Writer() *os.File {
	return s.writer
}
