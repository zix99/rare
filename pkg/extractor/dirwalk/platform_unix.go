//go:build unix

// Unix build tag is weird https://github.com/golang/go/issues/51572

package dirwalk

import (
	"os"
	"path/filepath"
	"syscall"
)

// return true if the dir is a different mount-point than its base
func isDifferentMount(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	statT, statOk := stat.Sys().(*syscall.Stat_t)
	if !statOk {
		return false
	}

	upDir := filepath.Dir(dir)
	upStat, err := os.Stat(upDir)
	if err != nil {
		return false
	}

	upT, upOk := upStat.Sys().(*syscall.Stat_t)
	if !upOk {
		return false
	}

	return statT.Dev != upT.Dev
}
