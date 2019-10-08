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

	inStatement := 0
	var sb strings.Builder
	runes := []rune(template)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\\' { // Escape
			i++
			sb.WriteRune(runes[i])
		} else if r == '{' {
			if inStatement == 0 { // starting a new token
				if sb.Len() > 0 {
					kb.stages = append(kb.stages, stageLiteral(sb.String()))
					sb.Reset()
				}
			} else {
				sb.WriteRune(r)
			}
			inStatement++
		} else if r == '}' && inStatement > 0 {
			inStatement--
			if inStatement == 0 {
				args := splitTokenizedArguments(sb.String())
				if len(args) == 1 { // Simple variable keyword like "{1}"
					kb.stages = append(kb.stages, stageSimpleVariable(args[0]))
				} else { // Complex function like "{add 1 2}"
					f := s.functions[args[0]]
					if f != nil {
						compiledArgs := make([]KeyBuilderStage, 0)
						for _, arg := range args[1:] {
							compiled := s.Compile(arg).joinStages()
							compiledArgs = append(compiledArgs, compiled)
						}
						kb.stages = append(kb.stages, f(compiledArgs))
					} else {
						kb.stages = append(kb.stages, stageError(fmt.Sprintf("Err:%s", args[0])))
					}
				}

				sb.Reset()
			} else {
				sb.WriteRune(r)
			}
		} else {
			sb.WriteRune(r)
		}
	}

	if sb.Len() > 0 && inStatement == 0 {
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

func (s *CompiledKeyBuilder) joinStages() KeyBuilderStage {
	if len(s.stages) == 0 {
		return stageLiteral("")
	}
	if len(s.stages) == 1 {
		return s.stages[0]
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		var sb strings.Builder
		for _, stage := range s.stages {
			sb.WriteString(stage(context))
		}
		return sb.String()
	})
}
