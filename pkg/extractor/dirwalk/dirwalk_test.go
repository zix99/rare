package dirwalk

import (
	"os"
	"path/filepath"
)

// test utils
var projectRoot = "./"

func init() {
	projectRoot = findProjectRoot()
	os.Chdir(projectRoot)
}

func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for !fileExists(filepath.Join(wd, "go.mod")) {
		up := filepath.Dir(wd)
		if wd == up {
			panic("unable to find wd")
		}
		wd = up
	}
	return wd
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func collectChan(c <-chan string) []string {
	ret := make([]string, 0)
	for s := range c {
		ret = append(ret, s)
	}
	return ret
}
