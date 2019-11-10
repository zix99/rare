package readahead

import (
	"bytes"
	"io"
)

/*
Buffered read-ahead similar to Scanner, except it will leave the large-buffers in place
(rather than shifting them) so that a given slice is good for the duration of its life

This allows a slice reference to be passed around without worrying that the underlying data will change
which limits the amount the data needs to be copied around

Initial benchmarks shows a 8% savings over Scanner
*/

type LineScanner interface {
	Scan() bool
	Bytes() []byte
}

type ReadAhead struct {
	r         io.Reader
	maxBufLen int

	buf    []byte
	offset int
	eof    bool

	token []byte
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func maxi(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func New(reader io.Reader, maxBufLen int) *ReadAhead {
	if maxBufLen <= 0 {
		panic("Buf length must be > 0")
	}
	return &ReadAhead{
		r:         reader,
		maxBufLen: maxBufLen,
	}
}

func (s *ReadAhead) Scan() bool {
	const averageLineLen = 512

	if !s.eof && (s.buf == nil || s.offset > len(s.buf)-averageLineLen) {
		// Need some more data in the buffer, but want to keep the existing one in-tact
		// This will lead to a small amount of duplication by copying the remaining data
		// in the old buf to the new buf
		oldbuf := s.buf
		s.buf = make([]byte, maxi(s.maxBufLen, len(oldbuf)-s.offset+s.maxBufLen/2))
		startRead := 0

		if oldbuf != nil && s.offset < len(oldbuf) {
			// Copy end of old buf over, and populate the rest of the buffer from reader
			startRead = len(oldbuf) - s.offset
			copy(s.buf, oldbuf[s.offset:])
		} else {
			// Brand new data!
			s.offset = 0
		}

		// Fill buffer
		for startRead < len(s.buf) {
			n, err := s.r.Read(s.buf[startRead:])
			startRead += n
			if err != nil {
				s.eof = true
				break
			}
		}

		s.buf = s.buf[:startRead]
		s.offset = 0
	}

	for {
		relIndex := bytes.IndexByte(s.buf[s.offset:], '\n')

		if relIndex >= 0 {
			start := s.offset
			s.offset += relIndex + 1
			s.token = dropCR(s.buf[start : start+relIndex])
			return true
		}

		// No new line, so either:
		// A) There's not enough room in the buffer
		// B) We're towards the end of buf, and need more data (Though hopefully we pre-empted above)
		if relIndex < 0 {
			// Eof, so the rest of it will count as a line
			if s.eof && s.offset < len(s.buf) {
				ret := s.buf[s.offset:]
				s.offset = len(s.buf)
				s.token = ret
				return true
			} else if !s.eof {
				// Not enough in buffer to find next new-line.. need to fill until finding
				oldbuf := s.buf
				s.buf = make([]byte, len(s.buf)*2)
				copy(s.buf, oldbuf)

				n, err := s.r.Read(s.buf[len(oldbuf):])
				s.buf = s.buf[:len(oldbuf)+n]

				if err != nil {
					s.eof = true
				}
			} else {
				s.token = nil
				return false
			}
		}
	}
}

func (s *ReadAhead) Bytes() []byte {
	return s.token
}

// ReadLine is shorthand for Scan() Token()
func (s *ReadAhead) ReadLine() []byte {
	if !s.Scan() {
		return nil
	}
	return s.token
}
