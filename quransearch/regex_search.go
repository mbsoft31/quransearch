package quransearch

import (
	"regexp"
	"time"
)

func (rx *RegexMethod) Search(text, p string, max int) []SearchMatch {
	start := time.Now()
	var matches = make([]SearchMatch, 0)
	re := regexp.MustCompile(p)
	for _, match := range re.FindAllStringIndex(text, max) {
		matches = append(matches, *NewSearchMatch(text, match[0], time.Since(start)))
	}
	return matches
}
