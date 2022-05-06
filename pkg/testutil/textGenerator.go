package testutil

import (
	"io"
	"math/rand"
)

type textGeneratingReader struct {
	maxChunk int
}

var _ io.Reader = &textGeneratingReader{}

var validText []byte = []byte("abcdefghijklmnopqrstuvwxyz\n\n")

// NewTextGenerator creates a io.reader that generates random alphaetical text separated by new-lines
// Will generate infinitely
func NewTextGenerator(maxReadSize int) io.Reader {
	return &textGeneratingReader{
		maxChunk: maxReadSize,
	}
}

func (s *textGeneratingReader) Read(buf []byte) (int, error) {
	size := len(buf)
	if size > s.maxChunk {
		size = s.maxChunk
	}

	for i := 0; i < size; i++ {
		buf[i] = validText[rand.Intn(len(validText))]
	}

	return size, nil
}
