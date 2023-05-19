package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/slicepool"
	"rare/pkg/stringSplitter"
	"strings"
)

// helper context to allow evaluating a sub-expression within a new context
type subContext struct {
	parent KeyBuilderContext
	vals   [2]string
}

var _ KeyBuilderContext = &subContext{}
var subContextPool = slicepool.NewObjectPool[subContext](5)

func (s *subContext) GetMatch(idx int) string {
	if idx < len(s.vals) {
		return s.vals[idx]
	}
	return ""
}

func (s *subContext) GetKey(k string) string {
	return s.parent.GetKey(k)
}

func (s *subContext) Eval(stage KeyBuilderStage, v0, v1 string) string {
	s.vals[0] = v0
	s.vals[1] = v1

	return stage(s)
}

// {@split <string> "delim"}
func kfArraySplit(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	byVal := EvalStageIndexOrDefault(args, 1, " ")
	if len(byVal) == 0 {
		return stageArgError(ErrEmpty, 1)
	}

	return func(context KeyBuilderContext) string {
		return arrayOperator(
			args[0](context),
			byVal,
			ArraySeparatorString,
			arrayOperatorNoopMapper,
		)
	}, nil
}

// {@join <array> "by"}
func kfArrayJoin(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 1, 2) {
		return stageErrArgRange(args, "1-2")
	}

	delim := EvalStageIndexOrDefault(args, 1, " ")
	return func(context KeyBuilderContext) string {
		return arrayOperator(
			args[0](context),
			ArraySeparatorString,
			delim,
			arrayOperatorNoopMapper,
		)
	}, nil
}

// {@select <array> "index"}
func kfArraySelect(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	index, indexOk := EvalStageInt(args[1])
	if !indexOk {
		return stageArgError(ErrNum, 1)
	}

	return func(context KeyBuilderContext) string {
		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		searchIndex := index
		if searchIndex < 0 {
			searchIndex += strings.Count(splitter.S, splitter.Delim) + 1
		}

		for i := 0; !splitter.Done(); i++ {
			val := splitter.Next()
			if i == searchIndex {
				return val
			}
		}

		return ""
	}, nil
}

// {@map <arr> <mapFunc>}
func kfArrayMap(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	return func(context KeyBuilderContext) string {
		mapperContext := subContextPool.Get()
		defer subContextPool.Return(mapperContext)

		*mapperContext = subContext{
			parent: context,
		}

		return arrayOperator(
			args[0](context),
			ArraySeparatorString,
			ArraySeparatorString,
			func(s string) string {
				return mapperContext.Eval(args[1], s, "")
			},
		)
	}, nil
}

// {@reduce <arr> <reducer> [initial=""]}
func kfArrayReduce(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 2, 3) {
		return stageErrArgRange(args, "2-3")
	}

	initial := EvalStageIndexOrDefault(args, 2, "")

	return func(context KeyBuilderContext) string {
		mapperContext := subContextPool.Get()
		defer subContextPool.Return(mapperContext)
		*mapperContext = subContext{
			parent: context,
		}

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		var memo string
		if initial == "" {
			memo = splitter.Next()
		} else {
			memo = initial
		}

		for !splitter.Done() {
			memo = mapperContext.Eval(args[1], memo, splitter.Next())
		}

		return memo
	}, nil
}

// {@slice <arr> start len}
func kfArraySlice(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if !isArgCountBetween(args, 2, 3) {
		return stageErrArgRange(args, "2-3")
	}

	sliceStart, ok := EvalStageInt(args[1])
	if !ok {
		return stageArgError(ErrConst, 1)
	}

	sliceLen, sliceLenOk := EvalArgInt(args, 2, -1)
	if !sliceLenOk {
		return stageArgError(ErrConst, 2)
	}

	return func(context KeyBuilderContext) string {
		var ret strings.Builder

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		realStart := sliceStart
		if realStart < 0 { // Negative start index starts from end
			realStart += strings.Count(splitter.S, ArraySeparatorString) + 1
		}

		for i := 0; (sliceLen < 0 || i < realStart+sliceLen) && !splitter.Done(); i++ {
			val := splitter.Next()
			if i >= realStart {
				if i > realStart {
					ret.WriteString(ArraySeparatorString)
				}
				ret.WriteString(val)
			}
		}

		return ret.String()
	}, nil
}

// {@filter <arr> <truthy-statement>}
func kfArrayFilter(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	return func(context KeyBuilderContext) string {
		var sb strings.Builder

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		sub := subContextPool.Get()
		defer subContextPool.Return(sub)
		*sub = subContext{
			parent: context,
		}

		needSep := false
		for !splitter.Done() {
			item := splitter.Next()

			if Truthy(sub.Eval(args[1], item, "")) {
				if needSep {
					sb.WriteRune(ArraySeparator)
				}
				sb.WriteString(item)
				needSep = true
			}
		}

		return sb.String()
	}, nil
}

func kfArrayIn(args []KeyBuilderStage) (KeyBuilderStage, error) {
	if len(args) != 2 {
		return stageErrArgCount(args, 2)
	}

	matchString, hasMatchString := EvalStaticStage(args[1])
	if !hasMatchString {
		return stageArgError(ErrConst, 1)
	}

	matchSet := make(map[string]struct{})
	for _, val := range strings.Split(matchString, ArraySeparatorString) {
		matchSet[val] = struct{}{}
	}

	return func(context KeyBuilderContext) string {
		val := args[0](context)
		if _, ok := matchSet[val]; ok {
			return TruthyVal
		}
		return FalsyVal
	}, nil
}

var arrayOperatorNoopMapper = func(s string) string { return s }

func arrayOperator(arr string, delim, joiner string, mapper func(o string) string) string {
	if arr == "" {
		return mapper(arr)
	}

	splitter := stringSplitter.Splitter{
		S:     arr,
		Delim: delim,
	}

	var ret strings.Builder
	ret.Grow(len(arr))

	// Map
	ret.WriteString(mapper(splitter.Next()))
	for !splitter.Done() {
		ret.WriteString(joiner)
		ret.WriteString(mapper(splitter.Next()))
	}

	return ret.String()
}
