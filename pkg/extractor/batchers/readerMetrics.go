package batchers

import "io"

type readerMetrics struct {
	r         io.Reader
	readBytes uint64
}

var _ io.Reader = &readerMetrics{}

func newReaderMetrics(r io.Reader) *readerMetrics {
	return &readerMetrics{r, 0}
}

func (s *readerMetrics) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)
	s.readBytes += uint64(n)
	return
}

func (s *readerMetrics) CountReset() (ret uint64) {
	ret = s.readBytes
	s.readBytes = 0
	return
}
