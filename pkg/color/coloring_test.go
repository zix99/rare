package color

import (
	"fmt"
	"testing"
)

func BenchmarkColorReplacer(b *testing.B) {
	s := "This is a test"
	groups := []int{5, 7, 8, 9}

	var out string
	for n := 0; n < b.N; n++ {
		out = WrapIndices(s, groups)
	}

	fmt.Println(out)
}
