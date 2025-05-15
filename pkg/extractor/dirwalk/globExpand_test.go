package dirwalk

import (
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
}

func TestGlobDirExclude(t *testing.T) {
	walk := Walker{
		ExcludeDir: []string{"cm?"},
	}
	p := "*/*.go"
	files := collectChan(walk.Walk(p))

	assertNoneContains(t, files, "cmd")
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
}

func TestRecurseInclude(t *testing.T) {
	walk := Walker{
		Recursive: true,
		Include:   []string{"*.sh", "*.go"},
	}

	files := collectChan(walk.Walk("docs/"))

	assert.Greater(t, len(files), 1)
	assertNoneContains(t, files, ".md")
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

func TestRecurseWithSymDir(t *testing.T) {
	if !fileExists("testwalk") {
		t.Skip("Needs testwalk dir")
	}

	walk := Walker{
		Recursive: true,
	}

	files := collectChan(walk.Walk("testwalk/"))
	assertNoneContains(t, files, "syminner")

	walk.FollowSymLinks = true
	files = collectChan(walk.Walk("testwalk/"))
	assert.Contains(t, files, "testwalk/syminner/infile")
}

func TestRecurseDoesntIdentifyDirAsFile(t *testing.T) {
	if !fileExists("testwalk") {
		t.Skip("Needs testwalk dir")
	}

	walk := Walker{
		Recursive:      true,
		FollowSymLinks: false,
		ListSymLinks:   true,
	}

	files := collectChan(walk.Walk("testwalk/"))

	assertNoneContains(t, files, "syminner")
}

func assertNoneContains(t *testing.T, set []string, contains string) {
	t.Helper()
	for _, item := range set {
		assert.NotContains(t, item, contains)
	}
}
