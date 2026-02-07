package helpers

import (
	"github.com/urfave/cli/v2"
)

const (
	ExitCodeNoData       = 1
	ExitCodeInvalidUsage = 2
	ExitCodeReadError    = 3
	ExitCodeOutputError  = 4
	ExitCodeSigInt       = 128 + 2 // 2 is SIGINT
)

type (
	BatcherErrors interface {
		ReadErrors() int
	}
	ExtractorSummary interface {
		MatchedLines() uint64
	}
	AggregationErrors interface {
		ParseErrors() uint64
	}
)

func DetermineErrorState(b BatcherErrors, e ExtractorSummary, agg AggregationErrors) error {
	if b.ReadErrors() > 0 {
		return cli.Exit("Read errors", ExitCodeReadError)
	}
	if agg != nil && agg.ParseErrors() > 0 {
		return cli.Exit("Parse errors", ExitCodeInvalidUsage)
	}
	if e.MatchedLines() == 0 {
		return cli.Exit("", ExitCodeNoData)
	}
	return nil
}

func DetermineErrorState2(interrupt bool, b BatcherErrors, e ExtractorSummary, agg AggregationErrors) error {
	if interrupt {
		return cli.Exit("", ExitCodeSigInt)
	}
	if b.ReadErrors() > 0 {
		return cli.Exit("Read errors", ExitCodeReadError)
	}
	if agg != nil && agg.ParseErrors() > 0 {
		return cli.Exit("Parse errors", ExitCodeInvalidUsage)
	}
	if e.MatchedLines() == 0 {
		return cli.Exit("", ExitCodeNoData)
	}
	return nil
}
