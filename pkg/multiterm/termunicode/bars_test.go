package termunicode

import (
	"os"
	"testing"
)

func TestWriteBar(t *testing.T) {
	BarWrite(os.Stdout, 1, 8, 1)
	BarWrite(os.Stdout, 10, 8, 1)
	BarString(5, 8, 1)
}
