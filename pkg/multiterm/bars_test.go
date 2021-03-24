package multiterm

import (
	"os"
	"testing"
)

func TestWriteBar(t *testing.T) {
	BarWrite(os.Stdout, 1, 8, 1)
}
