package funcfile

import (
	"rare/pkg/expressions"
	"rare/pkg/expressions/stdlib"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testContext = expressions.KeyBuilderContextArray{
	Elements: []string{"ab", "cd", "123"},
	Keys: map[string]string{
		"test": "testval",
	},
}

func TestCustomFunc(t *testing.T) {
	k := stdlib.NewStdKeyBuilder()
	_, err := createAndAddFunc(k, "double", "{sumi {0} {0} 4}")
	assert.NoError(t, err)

	kb, err := k.Compile("kd: {double {2}}")
	assert.Nil(t, err)
	val := kb.BuildKey(&testContext)
	assert.Equal(t, val, "kd: 250")
}

func TestCustomEdgeCases(t *testing.T) {
	k := stdlib.NewStdKeyBuilder()
	_, err := createAndAddFunc(k, "err", "{unclosed func")
	assert.Error(t, err)

	_, err = createAndAddFunc(k, "doublesrc", "{test} {sumi {0} {0}} missing: {1}")
	assert.NoError(t, err)
	kb, _ := k.Compile("{doublesrc 5}")
	assert.Equal(t, "testval 10 missing: ", kb.BuildKey(&testContext))
}

// BenchmarkCustomFunc-4   	 4767489	       244.2 ns/op	      16 B/op	       2 allocs/op
// BenchmarkCustomFunc-4   	 7563214	       160.4 ns/op	       3 B/op	       1 allocs/op
func BenchmarkCustomFunc(b *testing.B) {
	k := stdlib.NewStdKeyBuilder()
	createAndAddFunc(k, "double", "{sumi {0} {0} 4}")

	kb, err := k.Compile("{double {2}}")
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		kb.BuildKey(&testContext)
	}
}

func TestDeferredResolve(t *testing.T) {
	k := stdlib.NewStdKeyBuilderEx(false) // disable optimization because static analysis will trigger panic
	k.Func("panic", func(kbs []expressions.KeyBuilderStage) (expressions.KeyBuilderStage, error) {
		return func(kbc expressions.KeyBuilderContext) string {
			panic("not supposed to get here")
		}, nil
	})

	_, ferr := createAndAddFunc(k, "panicmissing", "{coalesce {0} {panic {0}}}")
	assert.NoError(t, ferr)

	kb, err := k.Compile("{panicmissing {1}}")
	assert.Nil(t, err)
	assert.Equal(t, "cd", kb.BuildKey(&testContext))

	assert.PanicsWithValue(t, "not supposed to get here", func() {
		kb.BuildKey(&expressions.KeyBuilderContextArray{
			Elements: []string{},
		})
	})
}
