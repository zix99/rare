package helpers

import (
	"testing"
)

func TestDetermineErrorState(t *testing.T) {
	// reader := testutil.NewTextGenerator(100)
	// b := batchers.OpenReaderToChan("test", reader, 1, 1024)
	// ext, _ := extractor.New(b.BatchChan(), &extractor.Config{
	// 	Regex:   ".*",
	// 	Extract: "{0}",
	// })
	// agg := aggregation.NewCounter()

	// for batch := range ext.ReadChan() {
	// 	for _, item := range batch {
	// 		agg.Sample(item.Extracted)
	// 	}
	// 	reader.Close() // Close soon after reading /some/ data
	// }

	// assert.NoError(t, DetermineErrorState(b, ext, agg))
}
