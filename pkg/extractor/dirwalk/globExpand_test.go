package dirwalk

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Though a little weird, the best we can do
is use facts about this repo to test the directory
navigation structure

We're using things that shouldn't change, or at least not
frequently, so should be safe.. but may need to update
these on refactors
*/

func TestDefaultOnDir(t *testing.T) {
	walk := Walker{}
	files := collectChan(walk.Walk("./"))

	// Returns dir even though not file
	assert.ElementsMatch(t, []string{"./"}, files)
}

func TestDefaultOnMissing(t *testing.T) {
	walk := Walker{}
	p := "not-exist"
	files := collectChan(walk.Walk(p))

	// Returns dir even though not file
	assert.ElementsMatch(t, []string{p}, files)
}

func TestDefaultOnGlob(t *testing.T) {
	walk := Walker{}
	p := "go.*"
	files := collectChan(walk.Walk(p))

	assert.ElementsMatch(t, []string{"go.mod", "go.sum"}, files)
}

func TestGlobNoExist(t *testing.T) {
	walk := Walker{}
	files := collectChan(walk.Walk("no-exist*"))
	assert.ElementsMatch(t, []string{"no-exist*"}, files)
}

func TestGlobInclude(t *testing.T) {
	walk := Walker{
		Include: []string{"*.mod"},
	}
	p := "go.*"
	files := collectChan(walk.Walk(p))

	assert.ElementsMatch(t, []string{"go.mod"}, files)
	assert.Len(t, files, 1)
}

func TestGlobExclude(t *testing.T) {
	walk := Walker{
		Exclude: []string{"*.sum"},
	}
	p := "go.*"
	files := collectChan(walk.Walk(p))

	assert.ElementsMatch(t, []string{"go.mod"}, files)
	assert.Len(t, files, 1)
	assert.Equal(t, uint64(1), walk.ExcludedCount())
}

func TestGlobDirExclude(t *testing.T) {
	walk := Walker{
		ExcludeDir: []string{"cm?"},
	}
	p := "*/*.go"
	files := collectChan(walk.Walk(p))

	assertNoneContains(t, files, "cmd")
	assert.Greater(t, walk.ExcludedCount(), uint64(1))
}

func TestRecurse(t *testing.T) {
	walk := Walker{
		Recursive: true,
	}
	p := "pkg/testutil"
	files := collectChan(walk.Walk(p))

	assert.Greater(t, len(files), 2)
}

func TestRecurseMissing(t *testing.T) {
	walk := Walker{
		Recursive: true,
	}
	files := collectChan(walk.Walk("missing/"))
	assert.Equal(t, []string{"missing/"}, files)
}

func TestRecurseNotDir(t *testing.T) {
	walk := Walker{
		Recursive: true,
	}
	files := collectChan(walk.Walk("go.mod"))
	assert.Equal(t, []string{"go.mod"}, files)
}

func TestRecurseExclude(t *testing.T) {
	walk := Walker{
		Recursive: true,
		Exclude:   []string{"*.sh", "*.go"},
	}

	files := collectChan(walk.Walk("docs/"))

	assert.Greater(t, len(files), 2)
	assertNoneContains(t, files, ".go")
	assertNoneContains(t, files, ".sh")
	assert.Greater(t, walk.ExcludedCount(), uint64(1))
}

func TestRecurseInclude(t *testing.T) {
	walk := Walker{
		Recursive: true,
		Include:   []string{"*.sh", "*.go"},
	}

	files := collectChan(walk.Walk("docs/"))

	assert.Greater(t, len(files), 1)
	assertNoneContains(t, files, ".md")
	assert.Greater(t, walk.ExcludedCount(), uint64(1))
}

func TestRecurseExcludeDir(t *testing.T) {
	walk := Walker{
		Recursive:  true,
		ExcludeDir: []string{"imag*", "usage"},
	}

	files := collectChan(walk.Walk("docs/"))

	assert.Greater(t, len(files), 1)
	assertNoneContains(t, files, "images")
	assertNoneContains(t, files, "usage")
	assert.Greater(t, walk.ExcludedCount(), uint64(1))
}

func TestRecurseWithSymFile(t *testing.T) {
	walk := Walker{
		Recursive:    true,
		ListSymLinks: true,
	}

	files := collectChan(walk.Walk("docs/"))
	assert.Contains(t, files, "docs/license.md")

	walk.ListSymLinks = false
	files = collectChan(walk.Walk("docs/"))
	assert.NotContains(t, files, "docs/license.md")
}

func TestRecursiveWithSymFileIgnore(t *testing.T) {
	walk := Walker{
		Recursive:    true,
		ListSymLinks: true,
		Exclude:      []string{"license*"},
	}

	files := collectChan(walk.Walk("docs/"))
	assert.NotContains(t, files, "docs/license.md")
	assert.Equal(t, uint64(1), walk.ExcludedCount())
}

func TestRecurseWithSymDir(t *testing.T) {
	p := setupTestDir(t)

	walk := Walker{
		Recursive: true,
	}

	files := collectChan(walk.Walk(p))
	assertNoneContains(t, files, "syminner")

	walk.FollowSymLinks = true
	files = collectChan(walk.Walk(p))
	assert.Contains(t, files, p+"/other/syminner/b")
}

func TestRecurseDoesntIdentifyDirAsFile(t *testing.T) {
	p := setupTestDir(t)

	walk := Walker{
		Recursive:      true,
		FollowSymLinks: false,
		ListSymLinks:   true,
	}

	files := collectChan(walk.Walk(p))

	assertNoneContains(t, files, "syminner")
}

func TestNoInfiniteRecursion(t *testing.T) {
	p := setupTestDir(t)
	os.Symlink("./", p+"/recursive")
	os.Symlink(p, p+"/recursive2")

	hadError := false
	walker := Walker{
		Recursive:       true,
		FollowSymLinks:  true,
		OnTraverseError: captureError(&hadError),
	}

	files := collectChan(walker.Walk(p))
	assert.True(t, hadError)
	assertNoneContains(t, files, "recursive")
	assertNoneContains(t, files, "recursive2")
	assert.Equal(t, uint64(0), walker.ExcludedCount())
}

func TestNoMountTraverseWithSymlink(t *testing.T) {
	p := setupTestDir(t)
	os.Symlink("/dev", p+"/dev")

	hadError := false
	walker := Walker{
		Recursive:       true,
		FollowSymLinks:  true,
		NoMountTraverse: true,
		OnTraverseError: captureError(&hadError),
	}

	files := collectChan(walker.Walk(p))
	assertNoneContains(t, files, "dev")
	assert.False(t, hadError)
}

func TestExcludeSymDir(t *testing.T) {
	p := setupTestDir(t)

	hadError := false
	walker := Walker{
		Recursive:       true,
		FollowSymLinks:  true,
		ExcludeDir:      []string{"syminner"},
		OnTraverseError: captureError(&hadError),
	}

	files := collectChan(walker.Walk(p))
	assert.False(t, hadError)
	assertNoneContains(t, files, "syminner")
	assert.Equal(t, uint64(1), walker.ExcludedCount())
}

func TestNoDoubleTraverseSymlink(t *testing.T) {
	p := setupTestDir(t)
	op := t.TempDir()
	os.WriteFile(op+"/opfile", []byte("hello"), 0644)
	os.Symlink(op, p+"/op1")
	os.Symlink(op, p+"/op2")

	hadError := false
	walker := Walker{
		Recursive:       true,
		FollowSymLinks:  true,
		OnTraverseError: captureError(&hadError),
	}

	files := collectChan(walker.Walk(p))
	assert.Equal(t, 1, countContains(files, "op1"))
	assert.Equal(t, 0, countContains(files, "op2"))
	assert.True(t, hadError)
	assert.Equal(t, uint64(0), walker.ExcludedCount())
}

func assertNoneContains(t *testing.T, set []string, contains string) {
	t.Helper()
	for _, item := range set {
		assert.NotContains(t, item, contains)
	}
}

func countContains(set []string, contains string) int {
	count := 0
	for _, item := range set {
		if strings.Contains(item, contains) {
			count++
		}
	}
	return count
}

func captureError(target *bool) func(err error) {
	return func(err error) {
		fmt.Println(err)
		*target = true
	}
}

/*
	 Sets up the following files in a temp dir to test more complex scenarios
		/a - "hello"
		/inner/b - "hello"
		/other/syminner -> /inner
		/other/symfile -> /inner/b
*/
func setupTestDir(t *testing.T) string {
	t.Helper()

	p := t.TempDir()

	os.WriteFile(p+"/a", []byte("hello"), 0644)

	os.Mkdir(p+"/inner", 0755)
	os.WriteFile(p+"/inner/b", []byte("hello"), 0644)

	os.Mkdir(p+"/other", 0755)
	os.Symlink(p+"/inner", p+"/other/syminner")
	os.Symlink(p+"/inner/b", p+"/other/symfile")

	return p
}
