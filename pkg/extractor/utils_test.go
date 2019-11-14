package extractor

func unbatchMatches(c <-chan []Match) []Match {
	matches := make([]Match, 0)
	for batch := range c {
		for _, item := range batch {
			matches = append(matches, item)
		}
	}
	return matches
}
