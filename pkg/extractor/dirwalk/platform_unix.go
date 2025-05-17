//go:build unix

// Unix build tag is weird https://github.com/golang/go/issues/51572

package dirwalk

import (
	"os"
	"syscall"
)

const FeatureMountTraversal = true

type DeviceId uint64

// Return the ID of the device the path is associated with
func getDeviceId(path string) DeviceId {
	stat, err := os.Lstat(path)
	if err != nil {
		return 0
	}

	statT, statOk := stat.Sys().(*syscall.Stat_t)
	if !statOk {
		return 0
	}

	return DeviceId(statT.Dev)
}
