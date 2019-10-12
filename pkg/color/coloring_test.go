package color

import (
	"fmt"
	"testing"
)

func BenchmarkColorReplacer(b *testing.B) {
	s := "This is a test"
	groups := []string{"is", "test"}

	var out string
	for n := 0; n < b.N; n++ {
		out = ColorCodeGroups(s, groups)
	}

	fmt.Println(out)
}
