package pathmatch

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncludeFilename(t *testing.T) {
	pm := &PathMatcher{
		Include:    MatchSet{"foo.txt"},
		Exclude:    MatchSet{"bar.txt"},
		ExcludeDir: MatchSet{},
	}
	assert.True(t, pm.IncludeFilename("foo.txt"))
	assert.False(t, pm.IncludeFilename("bar.txt"))
	assert.False(t, pm.IncludeFilename("baz.txt"))

	pm.Include = MatchSet{}
	assert.True(t, pm.IncludeFilename("baz.txt"))
}

func TestIncludeDirPath(t *testing.T) {
	pm := &PathMatcher{
		ExcludeDir: MatchSet{"skip"},
	}
	path := filepath.Join("/tmp", "skip", "file.txt")
	assert.False(t, pm.IncludeDirPath(path))
	assert.False(t, pm.IncludeDirPath("/skip/ok/file.txt"))
	assert.False(t, pm.IncludeDirPath("skip/ok/file.txt"))
	assert.True(t, pm.IncludeDirPath("/tmp/ok/file.txt"))
	assert.True(t, pm.IncludeDirPath("./tmp/ok/file.txt"))
	assert.True(t, pm.IncludeDirPath("/tmp/ok/skip"))

	pm.ExcludeDir = MatchSet{}
	assert.True(t, pm.IncludeDirPath(path))
}

func TestExcludeDirName(t *testing.T) {
	pm := &PathMatcher{
		ExcludeDir: MatchSet{"skip"},
	}
	assert.True(t, pm.ExcludeDirName("skip"))
	assert.False(t, pm.ExcludeDirName("ok"))
}

func TestIncludeFullPath(t *testing.T) {
	pm := &PathMatcher{
		Include:    MatchSet{"file.txt"},
		Exclude:    MatchSet{},
		ExcludeDir: MatchSet{"skip"},
	}
	path := filepath.Join("/tmp", "skip", "file.txt")
	assert.False(t, pm.IncludeFullPath(path))
	assert.True(t, pm.IncludeFullPath("/tmp/ok/file.txt"))
}
