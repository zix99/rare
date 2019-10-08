package extractor

import (
	"fmt"
	"strconv"
	"strings"
)

// KeyBuilderContext defines how to get information during run-time
type KeyBuilderContext interface {
	GetMatch(idx int) string
}

// KeyBuilderStage is a stage within the compiled builder
type KeyBuilderStage func(KeyBuilderContext) string

// KeyBuilderContextArray is a simple implementation of context with an array of elements
type KeyBuilderContextArray struct {
	Elements []string
}

func (s *KeyBuilderContextArray) GetMatch(idx int) string {
	if idx >= 0 && idx < len(s.Elements) {
		return s.Elements[idx]
	}
	return "<OOB>"
}

func stageLiteral(s string) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return s
	})
}

func stageSimpleVariable(s string) KeyBuilderStage {
	index, _ := strconv.Atoi(s)
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return context.GetMatch(index)
	})
}

func stageError(msg string) KeyBuilderStage {
	errMessage := fmt.Sprintf("<%s>", msg)
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		return errMessage
	})
}

// TODO: This can exist on compiled func
func stageJoin(stages []KeyBuilderStage) KeyBuilderStage {
	if len(stages) == 0 {
		return stageError("Empty")
	}
	if len(stages) == 1 {
		return stages[0]
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		var sb strings.Builder
		for _, stage := range stages {
			sb.WriteString(stage(context))
		}
		return sb.String()
	})
}
