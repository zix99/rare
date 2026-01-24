package helpers

import (
	"github.com/zix99/rare/pkg/logger"
	"github.com/zix99/rare/pkg/multiterm/termscaler"

	"github.com/urfave/cli/v2"
)

var ScaleFlag = &cli.StringFlag{
	Name:  "scale",
	Usage: "Defines data-scaling (linear, log10, log2)",
	Value: "linear",
}

func BuildScaler(scalerName string) (termscaler.Scaler, error) {
	if scalerName == "" {
		return termscaler.ScalerLinear, nil
	}
	return termscaler.ScalerByName(scalerName)
}

func BuildScalerOrFail(scalerName string) termscaler.Scaler {
	s, err := BuildScaler(scalerName)
	if err != nil {
		logger.Fatal(ExitCodeInvalidUsage, err)
	}
	return s
}
