package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}
