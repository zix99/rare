package iterwalk

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkFilepath-4   	   96153	     11961 ns/op	     681 B/op	      19 allocs/op
func BenchmarkFilepath(b *testing.B) {
	for range b.N {
		filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error { return err })
	}
}

// BenchmarkIter-4   	   86143	     14648 ns/op	     681 B/op	      19 allocs/op
func BenchmarkIter(b *testing.B) {
	for range b.N {
		WalkDir("./", func(path string, d fs.DirEntry, err error) error { return err })
	}
}

func TestWalkDir_File(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "testfile.txt")
	os.WriteFile(file, []byte("hello"), 0644)

	var walked []string
	err := WalkDir(file, func(path string, d fs.DirEntry, err error) error {
		walked = append(walked, path)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, []string{file}, walked)
}

func TestWalkDir_Dir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644)

	var walked []string
	err := WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		walked = append(walked, path)
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, len(walked), 3)
}

func TestWalkDir_SkipDir(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "subdir")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "file.txt"), []byte("x"), 0644)

	var walked []string
	err := WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		walked = append(walked, path)
		if d != nil && d.IsDir() && path == sub {
			return filepath.SkipDir
		}
		return nil
	})
	assert.NoError(t, err)
	assert.Contains(t, walked, sub)
	assert.NotContains(t, walked, filepath.Join(sub, "file.txt"))
}

func TestWalkDir_SkipDir_Large(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	dir := t.TempDir()
	sub := filepath.Join(dir, "bigsubdir")
	os.Mkdir(sub, 0755)
	// Create many files in the subdirectory
	for i := range readBatchSize + 1 {
		os.WriteFile(filepath.Join(sub, "file_"+strconv.Itoa(i)+".txt"), []byte("x"), 0644)
	}

	var walked []string
	err := WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d != nil && !d.IsDir() {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			walked = append(walked, path)
		}
		return nil
	})
	assert.NoError(t, err)
	assert.Len(t, walked, 0)

	// None of the files in the skipped subdir should be walked
	for i := 0; i < readBatchSize+1; i++ {
		assert.NotContains(t, walked, filepath.Join(sub, "file_"+strconv.Itoa(i)+".txt"))
	}
}

func TestWalkDir_Error(t *testing.T) {
	err := WalkDir("/nonexistent/path", func(path string, d fs.DirEntry, err error) error {
		return err
	})
	assert.Error(t, err)
}
