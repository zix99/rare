package stdlib

import (
	"rare/pkg/expressions"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Try to input-search for panics or other errors
// Doesn't check for correctness, just tries to break things
func TestAllFuncsNoPanic(t *testing.T) {
	// tests numerous combinations to assure code paths work and never panic
	if testing.Short() {
		t.SkipNow()
	}

	var argSets = [][]string{
		{"0", "5", "10", "20", "30", "40", "50", "60", "70", "100"},
		{"50", "40", "30", "20", "10", "0"},
		{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		{"", "", "", "", "", "", "", "", "", ""},
	}

	for name, fn := range StandardFunctions {
		t.Run(name, func(t *testing.T) {

			// run function with each argset and subset
			for _, argSet := range argSets {

				literalSet := make([]expressions.KeyBuilderStage, 0, len(argSet))
				contextSet := make([]expressions.KeyBuilderStage, 0, len(argSet))

				for i, arg := range argSet {
					// Compile and execute using a set of literals
					literalSet = append(literalSet, literal(arg))

					assert.NotPanics(t, func() {
						compiledFunc, err := fn(literalSet)

						// if compiled, try to invoke
						if err == nil {
							compiledFunc(mockContext())
						}
					})

					// try the same thing, but with contextual inputs
					argIdx := i
					contextSet = append(contextSet, func(kbc expressions.KeyBuilderContext) string { return kbc.GetMatch(argIdx) })
					ctx := &expressions.KeyBuilderContextArray{
						Elements: argSet[:i+1],
					}

					assert.NotPanics(t, func() {
						compiledFunc, err := fn(contextSet)

						// if compiled, try to invoke
						if err == nil {
							compiledFunc(ctx)
						}
					})

				}
			}

		})

	}

}
