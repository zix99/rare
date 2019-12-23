package fuzzy

func Distance(a, b []rune) int {
	alen := len(a)
	blen := len(b)
	col := make([]int, len(a)+1)

	for y := 1; y <= alen; y++ {
		col[y] = y
	}

	for x := 1; x <= blen; x++ {
		col[0] = x
		lastkey := x - 1
		for y := 1; y <= alen; y++ {
			oldkey := col[y]
			var incr int
			if a[y-1] != b[x-1] {
				incr = 1
			}
			col[y] = min3(col[y]+1, col[y-1]+1, lastkey+incr)
			lastkey = oldkey
		}
	}
	return col[alen]
}

func DistanceString(a, b string) int {
	return Distance([]rune(a), []rune(b))
}

func DistanceStringRatio(a, b string) float32 {
	dist := DistanceString(a, b)
	sum := len(a) + len(b)
	return float32(sum-dist*2) / float32(sum)
}
