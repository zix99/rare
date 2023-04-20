package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}

// todo: funcLiteral

// TODO: funcError
func stageError(err funcError) (KeyBuilderStage, error) {
	return func(ctx KeyBuilderContext) string {
		return err.expr
	}, err.err
}
