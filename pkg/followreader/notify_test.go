package followreader

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFileNotifyTail(t *testing.T) {
	af := CreateAppendingTempFile()
	defer af.Close()

	tail, err := NewNotify(af.Name(), false)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	assertSequentialReads(t, tail, 10)

	assert.NoError(t, tail.Close())
}

func TestTailNotifyFileAppendingExisting(t *testing.T) {
	af := CreateAppendingTempFile()

	tail, err := NewNotify(af.Name(), false)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	assertSequentialReads(t, tail, 10)

	// Re-open process
	af.Stop()
	af = CreateAppendingFromFile(af.Name())

	assertSequentialReads(t, tail, 10)

	af.Close()
	assert.NoError(t, tail.Close())
}

func TestTailNotifyFileRecreatedReopen(t *testing.T) {
	af := CreateAppendingTempFile()

	tail, err := NewNotify(af.Name(), true)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	assertSequentialReads(t, tail, 10)

	// Re-open process
	af.Stop()
	tail.Drain()
	af.Close() // Delete

	af = CreateAppendingFromFile(af.Name())

	assertSequentialReads(t, tail, 10)

	af.Close()
	assert.NoError(t, tail.Close())
}

func TestTailNotifyFileDeletedCloses(t *testing.T) {
	af := CreateAppendingTempFile()

	tail, err := NewNotify(af.Name(), false)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	assertSequentialReads(t, tail, 10)

	// Close and should delete
	af.Close()

	// Read until we receive an EOF
	gotEof := false
	buf := make([]byte, 100)
	for i := 0; i < 100; i++ {
		n, err := tail.Read(buf)
		fmt.Printf("Got data: %d\n", n)
		if err == io.EOF {
			gotEof = true
			break
		} else if err != nil {
			assert.Fail(t, "Non-eof error")
		}
	}

	assert.True(t, gotEof)
	assert.NoError(t, tail.Close())
}

func TestWatchingNonExistantFile(t *testing.T) {
	tp := path.Join(os.TempDir(), fmt.Sprintf("go-test-%d", rand.Int()))

	tail, err := NewNotify(tp, true)
	assert.NoError(t, err)

	af := CreateAppendingFromFile(tp)

	assertSequentialReads(t, tail, 10)

	af.Close()
	tail.Close()
}

func TestWatchingNonExistingFileFails(t *testing.T) {
	tp := path.Join(os.TempDir(), fmt.Sprintf("go-test-%d", rand.Int()))
	tail, err := NewNotify(tp, false)

	assert.Nil(t, tail)
	assert.Error(t, err)
}

func TestNonBlockingSignal(t *testing.T) {
	c := make(chan struct{}, 1)
	assert.Len(t, c, 0)
	writeSignalNonBlock(c)
	writeSignalNonBlock(c)
	assert.Len(t, c, 1)
	assert.NotNil(t, <-c)
}

func TestNotifyClosedReaderReturnsEOF(t *testing.T) {
	af := CreateAppendingTempFile()
	defer af.Close()

	tail, err := NewNotify(af.f.Name(), false)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	tail.Close()
	n, err := tail.Read(nil)
	assert.Zero(t, n)
	assert.ErrorIs(t, err, io.EOF)
}
