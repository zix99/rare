package expressions

import (
	"errors"
	"fmt"
	"strings"
)

// KeyBuilder builds the compiled keybuilder
type KeyBuilder struct {
	functions map[string]KeyBuilderFunction
}

// CompiledKeyBuilder represents the compiled key-builder
// can be considered thread-safe
type CompiledKeyBuilder struct {
	stages []KeyBuilderStage
}

// NewKeyBuilder creates a new KeyBuilder
func NewKeyBuilder() *KeyBuilder {
	kb := &KeyBuilder{
		functions: make(map[string]KeyBuilderFunction),
	}
	return kb
}

// Funcs appends a map of functions to be used by the parser
func (s *KeyBuilder) Funcs(funcs map[string]KeyBuilderFunction) {
	for k, f := range funcs {
		s.functions[k] = f
	}
}

// Compile builds a new key-builder
func (s *KeyBuilder) Compile(template string) (*CompiledKeyBuilder, error) {
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
							compiled, err := s.Compile(arg)
							if err != nil {
								return nil, err
							}
							compiledArgs = append(compiledArgs, compiled.joinStages())
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

	if inStatement != 0 {
		return nil, errors.New("non-terminated statement in expression")
	}

	if sb.Len() > 0 {
		kb.stages = append(kb.stages, stageLiteral(sb.String()))
	}

	return kb.optimize(), nil
}

func (s *CompiledKeyBuilder) optimize() *CompiledKeyBuilder {
	ret := &CompiledKeyBuilder{
		stages: make([]KeyBuilderStage, 0, len(s.stages)),
	}

	var sb strings.Builder
	for _, stage := range s.stages {
		if constVal, ok := EvalStaticStage(stage); ok {
			sb.WriteString(constVal)
		} else {
			if sb.Len() > 0 {
				ret.stages = append(ret.stages, stageLiteral(sb.String()))
				sb.Reset()
			}
			ret.stages = append(ret.stages, stage)
		}
	}

	if sb.Len() > 0 {
		ret.stages = append(ret.stages, stageLiteral(sb.String()))
	}

	return ret
}

func (s *CompiledKeyBuilder) BuildKey(context KeyBuilderContext) string {
	if len(s.stages) == 0 {
		return ""
	}
	if len(s.stages) == 1 {
		return s.stages[0](context)
	}

	var sb strings.Builder

	for _, stage := range s.stages {
		sb.WriteString(stage(context))
	}

	return sb.String()
}

func (s *CompiledKeyBuilder) StageCount() int {
	return len(s.stages)
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
