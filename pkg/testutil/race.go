//go:build race

package testutil

func init() {
	IsRaceMode = true
}
