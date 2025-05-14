package dirwalk

import (
	"os"
	"path/filepath"
	"rare/pkg/logger"
)

// Can be instantiated directly without newer and defaults will be
// vanilla walker
type Walker struct {
	Include    []string
	Exclude    []string
	ExcludeDir []string

	ListSymLinks   bool // Emit symlinks for non-directories
	FollowSymLinks bool // Recursively walk symlinks
	Recursive      bool // If asked to walk a path, will recurse
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
	if s.Recursive && s.isFollowableDir(p) {
		s.recurseWalk(c, p)
	} else {
		s.globExpand(c, p)
	}
}

func (s *Walker) recurseWalk(c chan<- string, p string) {
	filepath.WalkDir(p, func(walkPath string, info os.DirEntry, err error) error {
		switch {
		case err != nil: // error
			logger.Printf("Path error: %v", err)

		case info.IsDir() && isInMatchSet(s.ExcludeDir, info.Name()): // skipped dir
			return filepath.SkipDir

		case info.Type()&os.ModeSymlink != 0 && s.isFollowableDir(walkPath): // sym link dir
			// TODO: Prevent infinite recursion case
			s.recurseWalk(c, walkPath+string(filepath.Separator))

		case s.ListSymLinks && info.Type()&os.ModeSymlink != 0 && s.doesPathQualify(info.Name()): // sym link file
			c <- walkPath

		case info.Type().IsRegular() && s.doesPathQualify(info.Name()): // regular file
			c <- walkPath
		}
		return nil
	})
}

func (s *Walker) globExpand(c chan<- string, p string) {
	expanded, err := filepath.Glob(p)
	if err != nil {
		logger.Printf("Path error: %v", err)
	} else if len(expanded) > 0 {
		for _, item := range expanded {
			c <- item
		}
	} else {
		c <- p
	}
}

// If a given path is followable (dir or symlink to dir, as settings allow)
func (s *Walker) isFollowableDir(p string) bool {
	fi, err := os.Lstat(p)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return true
	}

	if s.FollowSymLinks && fi.Mode()&os.ModeSymlink != 0 {
		return s.isFollowableDir(p + string(filepath.Separator))
	}

	return false
}

// check path against includes/excludes
func (s *Walker) doesPathQualify(base string) bool {
	// Not in exclude list
	if isInMatchSet(s.Exclude, base) {
		return false
	}

	// If include list, assure in include list
	if len(s.Include) > 0 && !isInMatchSet(s.Include, base) {
		return false
	}

	return true
}

// Check if any of name match `filepath.Match` in matchSet
// for include/exclude logic
func isInMatchSet(matchSet []string, name string) bool {
	for _, match := range matchSet {
		if ok, _ := filepath.Match(match, name); ok {
			return true
		}
	}

	return false
}
