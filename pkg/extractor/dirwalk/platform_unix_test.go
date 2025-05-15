//go:build unix

package dirwalk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDifferentMount(t *testing.T) {
	assert.False(t, isDifferentMount("/etc"))
	assert.True(t, isDifferentMount("/proc"))
}
