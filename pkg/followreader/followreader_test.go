package followreader

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// reads data and asserts all reads successful
func assertSequentialReads(t *testing.T, tail FollowReader, reads int) {
	buf := make([]byte, 100)
	for i := 0; i < reads; i++ {
		n, err := tail.Read(buf)
		assert.NoError(t, err)
		assert.NotZero(t, n)
	}
}

// Helper to create a file and write random data at random intervals
type testAppendingFile struct {
	f    *os.File
	Line []byte

	stop chan<- bool
	wg   sync.WaitGroup
}

func CreateAppendingFromFile(filename string) *testAppendingFile {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return createAppendingFileEx(f)
}

func CreateAppendingTempFile() *testAppendingFile {
	f, err := ioutil.TempFile("", "go-test-")
	if err != nil {
		panic(err)
	}
	return createAppendingFileEx(f)
}

func createAppendingFileEx(f *os.File) *testAppendingFile {
	ret := &testAppendingFile{
		f:    f,
		Line: []byte("test file 123\n"),
		wg:   sync.WaitGroup{},
	}

	ret.startWriteRandomData(1 * time.Millisecond)

	return ret
}

func (s *testAppendingFile) Name() string {
	return s.f.Name()
}

func (s *testAppendingFile) startWriteRandomData(interval time.Duration) {
	stop := make(chan bool)
	s.stop = stop
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-stop:
				return
			case <-time.After(interval):
				s.f.Write(s.Line)
			}
		}
	}()
}

// Stop writing to the file
func (s *testAppendingFile) Stop() {
	if s.stop != nil {
		s.stop <- true
		s.wg.Wait()
		s.stop = nil
	}
}

// Close and stop writing to the file
func (s *testAppendingFile) Close() {
	s.Stop()
	err := s.f.Close()
	if err != nil {
		panic(err)
	}

	err = os.Remove(s.f.Name())
	if err != nil {
		panic(err)
	}
}
