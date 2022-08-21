package helpers

import (
	"rare/pkg/aggregation"
	"rare/pkg/extractor"
	"rare/pkg/extractor/batchers"

	"github.com/urfave/cli/v2"
)

const (
	ExitCodeNoData       = 1
	ExitCodeInvalidUsage = 2
)

func DetermineErrorState(b *batchers.Batcher, e *extractor.Extractor, agg aggregation.Aggregator) error {
	if b.ReadErrors() > 0 {
		return cli.Exit("Read errors", ExitCodeInvalidUsage)
	}
	if agg != nil && agg.ParseErrors() > 0 {
		return cli.Exit("Parse errors", ExitCodeInvalidUsage)
	}
	if e.MatchedLines() == 0 {
		return cli.Exit("", ExitCodeNoData)
	}
	return nil
}
