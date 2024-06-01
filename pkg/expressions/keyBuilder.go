package expressions

import (
	"fmt"
	"strings"
)

// KeyBuilder builds the compiled keybuilder
type KeyBuilder struct {
	functions    map[string]KeyBuilderFunction
	autoOptimize bool
}

// CompiledKeyBuilder represents the compiled key-builder
// can be considered thread-safe
type CompiledKeyBuilder struct {
	stages []KeyBuilderStage
}

// NewKeyBuilder creates a new KeyBuilder
func NewKeyBuilderEx(optimize bool) *KeyBuilder {
	kb := &KeyBuilder{
		functions:    make(map[string]KeyBuilderFunction),
		autoOptimize: optimize,
	}
	return kb
}

func NewKeyBuilder() *KeyBuilder {
	return NewKeyBuilderEx(true)
}

// Funcs appends a map of functions to be used by the parser
func (s *KeyBuilder) Funcs(funcs map[string]KeyBuilderFunction) {
	for k, f := range funcs {
		s.Func(k, f)
	}
}

// Funcs adds a single function used by the parser
func (s *KeyBuilder) Func(name string, f KeyBuilderFunction) {
	s.functions[name] = f
}

// Compile builds a new key-builder, returning error(s) on build issues
// if the CompiledKeyBuilder is not nil, then something is still useable (albeit may have problems)
func (s *KeyBuilder) Compile(template string) (*CompiledKeyBuilder, *CompilerErrors) {
	kb := &CompiledKeyBuilder{
		stages: make([]KeyBuilderStage, 0),
	}

	errs := CompilerErrors{
		Expression: template,
	}

	startStatement := 0
	inStatement := 0
	var sb strings.Builder
	runes := []rune(template)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '\\' { // Escape
			i++
			sb.WriteRune(unescape(runes[i]))
		} else if r == '{' {
			if inStatement == 0 { // starting a new token
				if sb.Len() > 0 {
					kb.stages = append(kb.stages, stageLiteral(sb.String()))
					sb.Reset()
				}
				startStatement = i
			} else {
				sb.WriteRune(r)
			}
			inStatement++
		} else if r == '}' && inStatement > 0 {
			inStatement--
			if inStatement == 0 {
				args := splitTokenizedArguments(sb.String())
				if len(args) == 0 {
					errs.add(ErrorEmptyStatement, string(runes[startStatement:i+1]), startStatement)
				} else if len(args) == 1 { // Simple variable keyword like "{1}"
					kb.stages = append(kb.stages, stageSimpleVariable(args[0]))
				} else { // Complex function like "{add 1 2}"
					f := s.functions[args[0]]
					if f != nil {
						compiledArgs := make([]KeyBuilderStage, 0, len(args)-1)
						for _, arg := range args[1:] {
							compiled, err := s.Compile(arg)
							if err != nil {
								errs.inherit(err, startStatement)
							}
							if compiled != nil {
								compiledArgs = append(compiledArgs, compiled.joinStages())
							}
						}
						stage, err := f(compiledArgs)
						if err != nil {
							errs.add(err, sb.String(), startStatement)
						}
						if stage != nil {
							kb.stages = append(kb.stages, stage)
						}
					} else {
						kb.stages = append(kb.stages, stageLiteral(fmt.Sprintf("<Err:%s>", args[0])))
						errs.add(ErrorMissingFunction, sb.String(), startStatement)
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
		errs.add(ErrorUnterminated, string(runes[startStatement:]), startStatement)
	}

	if sb.Len() > 0 {
		kb.stages = append(kb.stages, stageLiteral(sb.String()))
	}

	if s.autoOptimize {
		kb = kb.optimize()
	}

	if !errs.empty() {
		return kb, &errs
	}
	return kb, nil
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

func unescape(r rune) rune {
	switch r {
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	}
	return r
}
