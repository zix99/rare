package stdlib

import (
	. "rare/pkg/expressions" //lint:ignore ST1001 Legacy
	"rare/pkg/stringSplitter"
	"strings"
)

// helper context to allow evaluating a sub-expression within a new context
type subContext struct {
	parent KeyBuilderContext
	vals   [2]string
}

var _ KeyBuilderContext = &subContext{}

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
func kfArraySplit(args []KeyBuilderStage) KeyBuilderStage {
	if !isArgCountBetween(args, 1, 2) {
		return stageLiteral(ErrorArgCount)
	}

	byVal := EvalStageIndexOrDefault(args, 1, " ")
	if len(byVal) == 0 {
		return stageLiteral(ErrorEmpty)
	}

	return func(context KeyBuilderContext) string {
		return arrayOperator(
			args[0](context),
			byVal,
			ArraySeparatorString,
			arrayOperatorNoopMapper,
		)
	}
}

// {@join <array> "by"}
func kfArrayJoin(args []KeyBuilderStage) KeyBuilderStage {
	if !isArgCountBetween(args, 1, 2) {
		return stageLiteral(ErrorArgCount)
	}

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

// {@map <arr> <mapFunc>}
func kfArrayMap(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}

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

// {@reduce <arr> <reducer>}
func kfArrayReduce(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}

	return func(context KeyBuilderContext) string {
		mapperContext := subContext{
			parent: context,
		}

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		memo := splitter.Next()
		for !splitter.Done() {
			memo = mapperContext.Eval(args[1], memo, splitter.Next())
		}
		return memo
	}
}

// {@slice <arr> start len}
func kfArraySlice(args []KeyBuilderStage) KeyBuilderStage {
	if !isArgCountBetween(args, 2, 3) {
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
	}
}

// {@filter <arr> <truthy-statement>}
func kfArrayFilter(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageLiteral(ErrorArgCount)
	}

	return func(context KeyBuilderContext) string {
		var sb strings.Builder

		splitter := stringSplitter.Splitter{
			S:     args[0](context),
			Delim: ArraySeparatorString,
		}

		sub := subContext{
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