package funcfile

import (
	"github.com/zix99/rare/pkg/expressions"
	"github.com/zix99/rare/pkg/slicepool"
)

type lazySubContext struct {
	args []expressions.KeyBuilderStage
	sub  expressions.KeyBuilderContext
}

func (s *lazySubContext) GetMatch(idx int) string {
	if idx < 0 || idx >= len(s.args) {
		return ""
	}
	return s.args[idx](s.sub)
}

func (s *lazySubContext) GetKey(name string) string {
	return s.sub.GetKey(name)
}

func keyBuilderToFunction(stage *expressions.CompiledKeyBuilder) expressions.KeyBuilderFunction {
	return func(args []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
		ctxPool := slicepool.NewObjectPoolEx(5, func() *lazySubContext {
			return &lazySubContext{
				args: args,
			}
		})

		return func(kbc expressions.KeyBuilderContext) string {
			subCtx := ctxPool.Get()
			defer ctxPool.Return(subCtx)
			subCtx.sub = kbc

			return stage.BuildKey(subCtx)
		}, nil
	}
}
