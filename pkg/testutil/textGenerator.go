package testutil

import (
	"io"
)

type textGeneratingReader struct {
	maxChunk int
	closed   bool
}

var _ io.Reader = &textGeneratingReader{}

var validText []byte = []byte("abcdefghijklmnopqrstuvwxyz\n")

// NewTextGenerator creates a io.reader that generates random alphaetical text separated by new-lines
// Will generate infinitely
func NewTextGenerator(maxReadSize int) io.ReadCloser {
	return &textGeneratingReader{
		maxChunk: maxReadSize,
	}
}

func (s *textGeneratingReader) Read(buf []byte) (int, error) {
	if s.closed {
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

func (s *textGeneratingReader) Close() error {
	if s.closed {
		return io.EOF
	}
	s.closed = true
	return nil
}
