//go:build unix

package dirwalk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDifferentMount(t *testing.T) {
	rootDev := getDeviceId("/")
	assert.Equal(t, rootDev, getDeviceId("/etc"))
	assert.NotEqual(t, rootDev, getDeviceId("/proc"))
}

func BenchmarkGetDevId(b *testing.B) {
	for range b.N {
		getDeviceId("/proc")
	}
}
