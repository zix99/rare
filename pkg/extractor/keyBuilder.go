package extractor

import (
	"fmt"
	"strings"
)

// KeyBuilder builds the compiled keybuilder
type KeyBuilder struct {
	functions map[string]KeyBuilderFunction
}

// CompiledKeyBuilder represents the compiled key-builder
type CompiledKeyBuilder struct {
	stages []KeyBuilderStage
}

// NewKeyBuilder creates a new KeyBuilder
func NewKeyBuilder() *KeyBuilder {
	kb := &KeyBuilder{
		functions: make(map[string]KeyBuilderFunction),
	}
	kb.Funcs(defaultFunctions)
	return kb
}

// Funcs appends a map of functions to be used by the parser
func (s *KeyBuilder) Funcs(funcs map[string]KeyBuilderFunction) {
	for k, f := range funcs {
		s.functions[k] = f
	}
}

// Compile builds a new key-builder
func (s *KeyBuilder) Compile(template string) *CompiledKeyBuilder {
	kb := &CompiledKeyBuilder{
		stages: make([]KeyBuilderStage, 0),
	}

	inStatement := false
	var sb strings.Builder
	runes := []rune(template)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\\' { // Escape
			sb.WriteRune(runes[i+1])
			i++
		} else if r == '{' {
			kb.stages = append(kb.stages, stageLiteral(sb.String()))
			sb.Reset()
			inStatement = true
		} else if r == '}' && inStatement {
			keywords := strings.Split(sb.String(), " ")
			if len(keywords) == 1 { // Simple variable keyword like "{1}"
				kb.stages = append(kb.stages, stageSimpleVariable(keywords[0]))
			} else { // Complex function like "{add 1 2}"
				f := s.functions[keywords[0]]
				if f != nil {
					kb.stages = append(kb.stages, f(keywords[1:]))
				} else {
					kb.stages = append(kb.stages, stageError(fmt.Sprintf("Err:%s", keywords[0])))
				}
			}

			sb.Reset()
			inStatement = false
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 && !inStatement {
		kb.stages = append(kb.stages, stageLiteral(sb.String()))
	}

	return kb
}

func (s *CompiledKeyBuilder) BuildKey(context KeyBuilderContext) string {
	var sb strings.Builder

	for _, stage := range s.stages {
		sb.WriteString(stage(context))
	}

	return sb.String()
}
