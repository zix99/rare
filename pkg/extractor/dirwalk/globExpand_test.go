package dirwalk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobExpand(t *testing.T) {
	iter := GlobExpand([]string{"*"}, false)
	items := make([]string, 0)
	for ele := range iter {
		items = append(items, ele)
	}
	assert.Greater(t, len(items), 1)
}

func TestGlobExpandRecursive(t *testing.T) {
	iter := GlobExpand([]string{"../"}, true)
	items := make([]string, 0)
	for ele := range iter {
		items = append(items, ele)
	}
	assert.Greater(t, len(items), 10)
	fmt.Println(items)
}
