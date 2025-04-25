package helpers

import (
	"os"
	"os/signal"
	"rare/pkg/aggregation"
	"rare/pkg/extractor"
	"rare/pkg/logger"
	"sync"
	"time"
)

// RunAggregationLoop is a helper that takes care of output sync
// And the main async loops for you, it has two inputs (in addition to the extractor)
//
//	matchProcessor - to process a match
//	writeOutput - triggered after a delay, only if there's an update
//
// The two functions are guaranteed to never happen at the same time
func RunAggregationLoop(ext *extractor.Extractor, aggregator aggregation.Aggregator, writeOutput func()) {
	logger.DeferLogs()

	// Updater sync variables
	outputDone := make(chan bool)
	var outputMutex sync.Mutex

	// Updating loop
	go func() {
		for {
			select {
			case <-outputDone:
				return
			case <-time.After(100 * time.Millisecond):
				outputMutex.Lock()
				writeOutput()
				outputMutex.Unlock()
			}
		}
	}()

	// Processing data from extractor
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)
	reader := ext.ReadSimple()
PROCESSING_LOOP:
	for {
		select {
		case <-exitSignal:
			break PROCESSING_LOOP
		case matchBatch, more := <-reader:
			if !more {
				break PROCESSING_LOOP
			}
			outputMutex.Lock()
			for _, match := range matchBatch {
				aggregator.Sample(match)
			}
			outputMutex.Unlock()
		}
	}
	outputDone <- true

	writeOutput()
}
