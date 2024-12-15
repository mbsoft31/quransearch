package quransearch

import (
	"time"
)

// Search performs a brute-force search for the pattern in the given text
func (b *BruteForceMethod) Search(text, pattern string, max int) []SearchMatch {
	start := time.Now()
	matches := make([]SearchMatch, 0)

	m := len(text)
	n := len(pattern)

	if max < 0 {
		max = int(^uint(0) >> 1) // Max value for int
	}

	if n <= 0 || m <= 0 {
		return matches
	}

	for i := 0; i <= m-n && len(matches) < max; i++ {
		found := true
		for j := 0; j < n; j++ {
			if text[i+j] != pattern[j] {
				found = false
				break
			}
		}
		if found {
			elapsed := time.Since(start)
			match := NewSearchMatch(text, i, elapsed)
			matches = append(matches, *match)
			start = time.Now()
		}
	}

	return matches
}
