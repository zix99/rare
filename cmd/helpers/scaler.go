package helpers

import (
	"errors"
	"rare/pkg/logger"
	"rare/pkg/multiterm/termscaler"

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
	if scaler, ok := termscaler.ScalerByName(scalerName); ok {
		return scaler, nil
	}
	return termscaler.ScalerNull, errors.New("invalid scaler")
}

func BuildScalerOrFail(scalerName string) termscaler.Scaler {
	s, err := BuildScaler(scalerName)
	if err != nil {
		logger.Fatal(err)
	}
	return s
}
