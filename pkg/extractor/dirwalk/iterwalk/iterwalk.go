package iterwalk

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

/*
This is a fork of filepath.WalkDir that reimplements the os.ReadDir function
with one that opens and incrementally iterates on the paths without sorting.

For directories with many files, this both improves response time as well as
reducing the overall query time of the path.  For smaller directories it does
add some overhead, but it's negligible (microseconds)
*/

const readBatchSize = 1000

// This is an exact copy of `filepath.WalkDir`
func WalkDir(root string, fn fs.WalkDirFunc) error {
	info, err := os.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walkDir(root, fs.FileInfoToDirEntry(info), fn)
	}
	if err == filepath.SkipDir || err == filepath.SkipAll {
		return nil
	}
	return err
}

// Copy of filepath.walkDir, but switching from os.ReadDir to iterDir
func walkDir(path string, d fs.DirEntry, walkDirFn fs.WalkDirFunc) error {
	if err := walkDirFn(path, d, nil); err != nil || !d.IsDir() {
		if err == filepath.SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		// Second call, to report ReadDir error.
		err = walkDirFn(path, d, err)
		if err != nil {
			if err == filepath.SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}
	defer f.Close()

	for {
		dirs, err := f.ReadDir(readBatchSize)
		if err != nil && err != io.EOF {
			err = walkDirFn(path, d, err)
			if err != nil {
				if err == filepath.SkipDir && d.IsDir() {
					err = nil
				}
				return err
			}
		}

		if len(dirs) == 0 {
			break
		}

		for _, d1 := range dirs {
			path1 := filepath.Join(path, d1.Name())
			if err := walkDir(path1, d1, walkDirFn); err != nil {
				if err == filepath.SkipDir {
					break
				}
				return err
			}
		}
	}

	return nil
}
