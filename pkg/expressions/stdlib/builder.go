package stdlib

import "rare/pkg/expressions"

func NewStdKeyBuilder() *expressions.KeyBuilder {
	kb := expressions.NewKeyBuilder()
	kb.Funcs(StandardFunctions)
	return kb
}
