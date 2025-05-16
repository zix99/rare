package dirwalk

import (
	"fmt"
	"path/filepath"
)

// Set of filepath.Match globs with some helpers and error-checkers attached
type MatchSet []string

func NewMatchSet(patterns ...string) (MatchSet, error) {
	ms := MatchSet{}
	for _, p := range patterns {
		if err := ms.Add(p); err != nil {
			return nil, fmt.Errorf("error in '%s': %w", p, err)
		}
	}
	return ms, nil
}

func (s *MatchSet) Add(pattern string) error {
	_, err := filepath.Match(pattern, "")
	if err != nil {
		return err
	}
	*s = append(*s, pattern)
	return nil
}

func (s MatchSet) Matches(name string) bool {
	for _, pattern := range s {
		if ok, _ := filepath.Match(pattern, name); ok {
			return true
		}
	}
	return false
}
