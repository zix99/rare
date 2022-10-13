package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/stringSplitter"
	"strings"
)

type subContext struct {
	parent     KeyBuilderContext
	val0, val1 string
}

var _ KeyBuilderContext = &subContext{}

func (s *subContext) GetMatch(idx int) string {
	if idx == 0 {
		return s.val0
	}
	if idx == 1 {
		return s.val1
	}
	return s.parent.GetMatch(idx)
}

func (s *subContext) GetKey(k string) string {
	return s.parent.GetKey(k)
}

func (s *subContext) Eval(stage KeyBuilderStage, v0, v1 string) string {
	s.val0, s.val1 = v0, v1
	return stage(s)
}

// {split <string> "delim"}
func kfArraySplit(args []KeyBuilderStage) KeyBuilderStage {
	byVal := EvalStageOrDefault(args[1], " ")
	return func(context KeyBuilderContext) string {
		return arrayOperator(
			args[0](context),
			byVal,
			ArraySeparatorString,
			arrayOperatorNoopMapper,
		)
	}
}

// {join <array> "by"}
func kfArrayJoin(args []KeyBuilderStage) KeyBuilderStage {
	delim := EvalStageIndexOrDefault(args, 1, " ")
	return func(context KeyBuilderContext) string {
		return arrayOperator(
			args[0](context),
			ArraySeparatorString,
			delim,
			arrayOperatorNoopMapper,
		)
	}
}

// {map <arr> <mapFunc>}
func kfArrayMap(args []KeyBuilderStage) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		mapperContext := subContext{
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
	}
}

// {reduce <arr> <reducer>}
func kfArrayReduce(args []KeyBuilderStage) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		mapperContext := subContext{
			parent: context,
		}

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		mapperContext.val0 = splitter.Next()
		for !splitter.Done() {
			mapperContext.val1 = splitter.Next()
			mapperContext.val0 = args[1](&mapperContext)
		}
		return mapperContext.val0
	}
}

// {slice <arr> start len}
func kfArraySlice(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) < 2 {
		return stageLiteral(ErrorArgCount)
	}

	sliceStart := EvalStageInt(args[1], 0)
	var sliceLen int = -1
	if len(args) >= 3 {
		sliceLen = EvalStageInt(args[2], -1)
	}

	return func(context KeyBuilderContext) string {
		var ret strings.Builder

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		for i := 0; (sliceLen < 0 || i < sliceStart+sliceLen) && !splitter.Done(); i++ {
			val := splitter.Next()
			if i >= sliceStart {
				if i > sliceStart {
					ret.WriteString(ArraySeparatorString)
				}
				ret.WriteString(val)
			}
		}

		return ret.String()
	}
}

// {filter <arr> <truthy-statement>}
func kfArrayFilter(args []KeyBuilderStage) KeyBuilderStage {
	return func(context KeyBuilderContext) string {
		var sb strings.Builder

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		sub := subContext{
			parent: context,
		}

		for !splitter.Done() {
			item := splitter.Next()

			if Truthy(sub.Eval(args[1], item, "")) {
				if sb.Len() > 0 {
					sb.WriteRune(ArraySeparator)
				}
				sb.WriteString(item)
			}
		}

		return sb.String()
	}
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
