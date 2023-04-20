package stdlib

import (
	"path/filepath"
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
)

func kfPathManip(manipulator func(string) string) KeyBuilderFunction {
	return func(args []KeyBuilderStage) (KeyBuilderStage, error) {
		if len(args) != 1 {
			return stageError(ErrArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			return manipulator(args[0](context))
		}), nil
	}
}

var (
	kfPathBase = kfPathManip(filepath.Base)
	kfPathDir  = kfPathManip(filepath.Dir)
	kfPathExt  = kfPathManip(filepath.Ext)
)
