package multiterm

import (
	"os"
	"testing"
)

func TestWriteBar(t *testing.T) {
	WriteBar(os.Stdout, 1, 8, 1)
}
