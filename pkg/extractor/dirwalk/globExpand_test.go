package dirwalk

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirWalk(t *testing.T) {
	fmt.Println(os.Getwd())
	walk := Walker{}
	files := collectChan(walk.Walk("./dirtest"))
	assert.ElementsMatch(t, []string{}, files)
}

func collectChan(c <-chan string) []string {
	ret := make([]string, 0)
	for s := range c {
		ret = append(ret, s)
	}
	return ret
}

// func TestGlobExpand(t *testing.T) {
// 	iter := GlobExpand([]string{"*"}, false)
// 	items := make([]string, 0)
// 	for ele := range iter {
// 		items = append(items, ele)
// 	}
// 	assert.Greater(t, len(items), 1)
// }

// func TestGlobExpandRecursive(t *testing.T) {
// 	iter := GlobExpand([]string{"../"}, true)
// 	items := make([]string, 0)
// 	for ele := range iter {
// 		items = append(items, ele)
// 	}
// 	assert.Greater(t, len(items), 10)
// 	fmt.Println(items)
// }

// func TestGlobalExpandNoFile(t *testing.T) {
// 	iter := GlobExpand([]string{"does-not-exist"}, false)

// 	items := make([]string, 0)
// 	for ele := range iter {
// 		items = append(items, ele)
// 	}

// 	assert.Len(t, items, 1)
// }
