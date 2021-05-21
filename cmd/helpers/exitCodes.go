package helpers

import (
	"rare/pkg/aggregation"
	"rare/pkg/extractor"
	"rare/pkg/extractor/batchers"

	"github.com/urfave/cli"
)

const (
	ExitCodeNoData       = 1
	ExitCodeInvalidUsage = 2
)

func DetermineErrorState(b *batchers.Batcher, e *extractor.Extractor, agg aggregation.Aggregator) error {
	if b.ReadErrors() > 0 {
		return cli.NewExitError("Read errors", ExitCodeInvalidUsage)
	}
	if agg != nil && agg.ParseErrors() > 0 {
		return cli.NewExitError("Parse errors", ExitCodeInvalidUsage)
	}
	if e.MatchedLines() == 0 {
		return cli.NewExitError("", ExitCodeNoData)
	}
	return nil
}
