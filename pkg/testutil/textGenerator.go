package testutil

import (
	"io"
	"sync/atomic"
)

type textGeneratingReader struct {
	maxChunk int
	closed   int32
}

var _ io.Reader = &textGeneratingReader{}

var validText []byte = []byte("abcdefghijklmnopqrstuvwxyz\n")

// NewTextGenerator creates a io.reader that generates random alphaetical text separated by new-lines
// Will generate infinitely until closed
func NewTextGenerator(maxReadSize int) io.ReadCloser {
	return &textGeneratingReader{
		maxChunk: maxReadSize,
	}
}

func (s *textGeneratingReader) Read(buf []byte) (int, error) {
	if atomic.LoadInt32(&s.closed) > 0 {
		return 0, io.EOF
	}

	size := len(buf)
	if size > s.maxChunk {
		size = s.maxChunk
	}

	for i := 0; i < size; i += len(validText) {
		copy(buf[i:size], validText)
	}

	return size, nil
}

// Close, next Read() will return EOF (thread-safe, for testing)
func (s *textGeneratingReader) Close() error {
	if atomic.SwapInt32(&s.closed, 1) > 0 {
		return io.EOF
	}
	return nil
}
