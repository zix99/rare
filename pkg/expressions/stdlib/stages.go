package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func literal(s string) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		return s
	}
}

func stageLiteral(s string) (KeyBuilderStage, error) {
	return literal(s), nil
}
