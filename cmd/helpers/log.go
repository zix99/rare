package helpers

import (
	"log"
	"os"
)

var stderrLog = log.New(os.Stderr, "[Log] ", 0)
