package pathmatch

import "path/filepath"

type PathMatcher struct {
	Include    MatchSet
	Exclude    MatchSet
	ExcludeDir MatchSet
}

// Check only basename against includes/excludes
func (s *PathMatcher) IncludeFilename(basename string) bool {
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

// SLOW: Recursively check all sub-paths of path (assumes is a path to a file)
func (s *PathMatcher) IncludeDirPath(fullpath string) bool {
	if len(s.ExcludeDir) == 0 { //shortcut
		return true
	}

	cur := filepath.Dir(fullpath)
	for cur != "." && cur != string(filepath.Separator) {
		if s.ExcludeDir.Matches(filepath.Base(cur)) {
			return false
		}
		cur = filepath.Dir(cur)
	}

	return true
}

// Check if directory name is excluded in ExcludeDir
func (s *PathMatcher) ExcludeDirName(basename string) bool {
	return s.ExcludeDir.Matches(basename)
}

// SLOW: Check both filename and recurse the entire path for dir-excludes
func (s *PathMatcher) IncludeFullPath(fullpath string) bool {
	return s.IncludeFilename(filepath.Base(fullpath)) && s.IncludeDirPath(fullpath)
}
