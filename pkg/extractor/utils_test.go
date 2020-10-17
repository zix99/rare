package extractor

func unbatchMatches(c <-chan []Match) []Match {
	matches := make([]Match, 0)
	for batch := range c {
		matches = append(matches, batch...)
	}
	return matches
}
