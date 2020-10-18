package stdlib

import (
	"fmt"

	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}

func stageError(msg string) KeyBuilderStage {
	errMessage := fmt.Sprintf("<%s>", msg)
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return errMessage
	})
}
