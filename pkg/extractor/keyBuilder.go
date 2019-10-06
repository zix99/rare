package extractor

import (
	"strings"
)

// KeyBuilder builds the compiled keybuilder
type KeyBuilder struct {
	functions map[string]KeyBuilderFunction
}

// CompiledKeyBuilder represents the compiled key-builder
type CompiledKeyBuilder struct {
	stages  []KeyBuilderStage
	builder *KeyBuilder
}

func NewKeyBuilder() *KeyBuilder {
	kb := &KeyBuilder{
		functions: make(map[string]KeyBuilderFunction),
	}
	kb.Funcs(defaultFunctions)
	return kb
}

func (s *KeyBuilder) Funcs(funcs map[string]KeyBuilderFunction) {
	for k, f := range funcs {
		s.functions[k] = f
	}
}

// NewKeyBuilder builds a new key-builder
func (s *KeyBuilder) Compile(template string) *CompiledKeyBuilder {
	kb := &CompiledKeyBuilder{
		stages:  make([]KeyBuilderStage, 0),
		builder: s,
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
			if len(keywords) == 1 {
				kb.stages = append(kb.stages, stageSimpleVariable(keywords[0]))
			} else {
				f := s.functions[keywords[0]]
				kb.stages = append(kb.stages, f(keywords[1:]))
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
