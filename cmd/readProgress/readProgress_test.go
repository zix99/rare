package readProgress

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadProgress(t *testing.T) {
	SetSourceCount(5)
	IncSourceCount(1)

	assert.Equal(t, 6, sourceCount)

	StartFileReading("abc")
	assert.Contains(t, GetReadFileString(), "0/6")

	StopFileReading("abc")
	assert.Equal(t, 6, sourceCount)
	assert.Equal(t, 1, readCount)
	assert.Contains(t, GetReadFileString(), "1/6")
}
