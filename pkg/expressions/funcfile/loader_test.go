package funcfile

import (
	"rare/pkg/expressions/stdlib"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefinitions(t *testing.T) {
	data := `# a comment
	double {sumi {0} {0}}
	
	# Another comment
	quad this equals: \
	  {multi \
		{0} \
		{0} \
		#Mixed comment
		{0} {0}}
	`
	kb := stdlib.NewStdKeyBuilder()
	funcs, err := LoadDefinitions(kb, strings.NewReader(data))
	assert.NoError(t, err)
	assert.Len(t, funcs, 2)
	assert.Contains(t, funcs, "quad")
	assert.Contains(t, funcs, "double")

	val, err := kb.Compile("{quad 5}")
	assert.Nil(t, err)
	assert.Equal(t, "this equals: 625", val.BuildKey(nil))
}

func TestLoadDefinitionsErrs(t *testing.T) {
	data := `# a comment
	unterm unterm {
	nofunc

	`
	kb := stdlib.NewStdKeyBuilder()
	funcs, err := LoadDefinitions(kb, strings.NewReader(data))
	assert.NotNil(t, err)
	assert.Len(t, funcs, 0)
}

func TestLoadFile(t *testing.T) {
	kb := stdlib.NewStdKeyBuilder()
	funcs, err := LoadDefinitionsFile(kb, "example.funcfile")
	assert.NoError(t, err)
	assert.Len(t, funcs, 1)

	_, err = LoadDefinitionsFile(kb, "notfile.funcfile")
	assert.Error(t, err)
}
