package dirwalk

import (
	"fmt"
	"os"
	"path/filepath"
	"rare/pkg/extractor/dirwalk/iterwalk"
	"sync/atomic"
)

// Can be instantiated directly without newer and defaults will be
// vanilla walker
type Walker struct {
	Include    MatchSet
	Exclude    MatchSet
	ExcludeDir MatchSet

	ListSymLinks    bool // Emit symlinks for non-directories
	FollowSymLinks  bool // Recursively walk symlinks
	Recursive       bool // If asked to walk a path, will recurse
	NoMountTraverse bool // If true, don't traverse into another mount

	OnTraverseError func(error) // Called in separate goroutine

	total    atomic.Uint64
	excluded atomic.Uint64 // Files excluded due to include/exclude rules (not sym or mount rules)
	error    atomic.Uint64
}

type Metrics interface {
	TotalCount() uint64
	ExcludedCount() uint64
	ErrorCount() uint64
}

// Number of paths skipped because of rules (include, exclude, exludedir; NOT skip sym, mounts, etc)
func (s *Walker) ExcludedCount() uint64 {
	return s.excluded.Load()
}

func (s *Walker) TotalCount() uint64 {
	return s.total.Load()
}

func (s *Walker) ErrorCount() uint64 {
	return s.error.Load()
}

func (s *Walker) Walk(paths ...string) <-chan string {
	c := make(chan string, 10)

	go func() {
		for _, p := range paths {
			s.walk(c, p)
		}
		close(c)
	}()

	return c
}

func (s *Walker) walk(c chan<- string, p string) {
	if s.Recursive && isFollowableDir(p) {
		s.recurseWalk(c, p, map[string]string{p: p})
	} else {
		s.globExpand(c, p)
	}
}

func (s *Walker) recurseWalk(c chan<- string, p string, visited map[string]string) {
	var rootDevId DeviceId
	if s.NoMountTraverse {
		// getDeviceId (stat) is expensive
		rootDevId = getDeviceId(p)
	}

	iterwalk.WalkDir(p, func(walkPath string, info os.DirEntry, err error) error {
		switch {
		case err != nil: // error
			s.onError(fmt.Errorf("path error: %w", err))

		case info.IsDir() && s.ExcludeDir.Matches(info.Name()): // skipped dir
			s.excluded.Add(1)
			return filepath.SkipDir

		case info.IsDir() && s.NoMountTraverse && getDeviceId(walkPath) != rootDevId: // skipped mount
			return filepath.SkipDir

		case info.Type()&os.ModeSymlink != 0 && isFollowableDir(walkPath): // sym link dir
			// WalkDir won't navigate symlinks by default. This will traverse recursively
			if !s.FollowSymLinks {
				break
			}

			if s.ExcludeDir.Matches(info.Name()) {
				s.excluded.Add(1)
				break
			}

			real, err := filepath.EvalSymlinks(walkPath)
			if err != nil {
				s.onError(err)
			} else if s.NoMountTraverse && getDeviceId(real) != rootDevId {
				// skip
			} else if with, ok := visited[real]; ok {
				s.onError(fmt.Errorf("already traversed symlink %s in %s", walkPath, with))
			} else {
				visited[real] = walkPath
				s.recurseWalk(c, walkPath+string(filepath.Separator), visited)
			}

		case info.Type()&os.ModeSymlink != 0: // sym link file
			if !s.ListSymLinks {
				break
			}

			if !s.shouldIncludeFilename(info.Name()) {
				s.excluded.Add(1)
				break
			}

			c <- walkPath
			s.total.Add(1)

		case info.Type().IsRegular(): // regular file
			if !s.shouldIncludeFilename(info.Name()) {
				s.excluded.Add(1)
				break
			}
			c <- walkPath
			s.total.Add(1)
		}
		return nil
	})
}

// Uses glob expand, eg '*.txt'
func (s *Walker) globExpand(c chan<- string, p string) {
	expanded, err := filepath.Glob(p)
	if err != nil {
		s.onError(fmt.Errorf("path error: %w", err))
	} else if len(expanded) > 0 {
		for _, item := range expanded {
			if s.shouldIncludeFilename(filepath.Base(item)) && s.shouldIncludeDir(item) {
				c <- item
				s.total.Add(1)
			} else {
				s.excluded.Add(1)
			}
		}
	} else {
		c <- p
		s.total.Add(1)
	}
}

// check path against includes/excludes
func (s *Walker) shouldIncludeFilename(basename string) bool {
	// Not in exclude list
	if s.Exclude.Matches(basename) {
		return false
	}

	// If include list, assure in include list
	if len(s.Include) > 0 && !s.Include.Matches(basename) {
		return false
	}

	return true
}

// Takes in a full path eg abc/efg/filename
// Checks against ExcludeDir
func (s *Walker) shouldIncludeDir(fullpath string) bool {
	if len(s.ExcludeDir) == 0 { //shortcut
		return true
	}

	cur := filepath.Dir(fullpath)
	for cur != "." {
		if s.ExcludeDir.Matches(filepath.Base(cur)) {
			return false
		}
		cur = filepath.Dir(cur)
	}

	return true
}

func (s *Walker) onError(err error) {
	s.error.Add(1)
	if s.OnTraverseError != nil {
		s.OnTraverseError(err)
	}
}

// If a given path is followable (dir or symlink to dir)
func isFollowableDir(p string) bool {
	fi, err := os.Lstat(p)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return true
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		return isFollowableDir(p + string(filepath.Separator))
	}

	return false
}
