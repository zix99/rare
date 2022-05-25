package followreader

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFileNotifyTail(t *testing.T) {
	af := CreateAppendingTempFile()
	defer af.Close()

	tail, err := NewNotify(af.Name(), false)
	assert.NoError(t, err)
	assert.NotNil(t, tail)

	ret := make([]byte, 100)
	for i := 0; i < 10; i++ {
		n, err := tail.Read(ret)
		assert.NoError(t, err)
		assert.NotZero(t, n)
	}

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

// TODO:
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

	if !gotEof {
		assert.Fail(t, "Never received EOF")
	}

	assert.NoError(t, tail.Close())
}
