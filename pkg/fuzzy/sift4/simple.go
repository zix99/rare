package sift4

func max2i(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func DistanceString(s1, s2 string, maxOffset int) int {
	return Distance([]rune(s1), []rune(s2), maxOffset)
}

// Sift4 - simplest version
// online algorithm to compute the distance between two strings in O(n)
// maxOffset is the number of characters to search for matching letters
// https://siderite.dev/blog/super-fast-and-accurate-string-distance.html/#at929344888
func Distance(s1, s2 []rune, maxOffset int) int {
	if len(s1) == 0 {
		if len(s2) == 0 {
			return 0
		}
		return len(s2)
	}

	if len(s2) == 0 {
		return len(s1)
	}

	l1 := len(s1)
	l2 := len(s2)

	c1 := 0       //cursor for string 1
	c2 := 0       //cursor for string 2
	lcss := 0     //largest common subsequence
	local_cs := 0 //local common substring
	for {
		if c1 >= l1 || c2 >= l2 {
			break
		}

		if s1[c1] == s2[c2] {
			local_cs++
		} else {
			lcss += local_cs
			local_cs = 0
			if c1 != c2 {
				c1 = max2i(c1, c2) //using max to bypass the need for computer transpositions ('ab' vs 'ba')
				c2 = c1
			}
			for i := 0; i < maxOffset && (c1+i < l1 || c2+i < l2); i++ {
				if c1+i < l1 && c2 < l2 && s1[c1+i] == s2[c2] {
					c1 += i
					local_cs++
					break
				}
				if c1 < l1 && c2+i < l2 && s1[c1] == s2[c2+i] {
					c2 += i
					local_cs++
					break
				}
			}
		}
		c1++
		c2++
	}
	lcss += local_cs
	return max2i(l1, l2) - lcss
}
