package quransearch

import (
	"time"
	"unicode/utf8"
)

// Search searches for occurrences of the pattern in the given text
func (bm *BoyerMooreMethod) Search(text, pattern string, max int) []SearchMatch {
	bm.Pattern = pattern
	bm.PatternLength = utf8.RuneCountInString(pattern)
	bm.BadCharacter = bm.makeBadCharacterShifts()
	bm.GoodSuffix = bm.makeGoodSuffixShifts()

	var matches []SearchMatch
	if bm.PatternLength == 0 {
		return matches
	}

	if max <= 0 {
		max = int(^uint(0) >> 1) // Max value for int
	}

	start := time.Now()

	textLength := utf8.RuneCountInString(text)
	i := bm.PatternLength - 1

	for i < textLength && len(matches) < max {
		j := bm.PatternLength - 1
		for j >= 0 && bm.Pattern[j] == text[i] {
			i--
			j--
		}
		if j < 0 {
			match := NewSearchMatch(text, i+1, time.Since(start))
			matches = append(matches, *match)
			i += bm.PatternLength * 2
		} else {
			i += maxInt(bm.GoodSuffix[bm.PatternLength-1-j], bm.BadCharacter[text[i]])
		}
	}

	return matches
}

func (bm *BoyerMooreMethod) makeBadCharacterShifts() []int {
	const AlphabetSize = 'ي' + 1

	badCS := make([]int, AlphabetSize)
	for i := 0; i < AlphabetSize; i++ {
		badCS[i] = bm.PatternLength
	}
	for i := 0; i < bm.PatternLength-1; i++ {
		badCS[bm.Pattern[i]] = bm.PatternLength - 1 - i
	}
	return badCS
}

func (bm *BoyerMooreMethod) makeGoodSuffixShifts() []int {
	goodSuffix := make([]int, bm.PatternLength)
	lastPrefixPosition := bm.PatternLength

	for i := bm.PatternLength - 1; i >= 0; i-- {
		if bm.isPrefix(i + 1) {
			lastPrefixPosition = i + 1
		}
		goodSuffix[bm.PatternLength-1-i] = lastPrefixPosition - i + bm.PatternLength - 1
	}

	for i := 0; i < bm.PatternLength-1; i++ {
		slen := bm.suffixLength(i)
		goodSuffix[slen] = bm.PatternLength - 1 - i + slen
	}

	return goodSuffix
}

func (bm *BoyerMooreMethod) isPrefix(p int) bool {
	for i, j := p, 0; i < bm.PatternLength; i, j = i+1, j+1 {
		if bm.Pattern[i] != bm.Pattern[j] {
			return false
		}
	}
	return true
}

func (bm *BoyerMooreMethod) suffixLength(p int) int {
	length := 0
	for i, j := p, bm.PatternLength-1; i >= 0 && bm.Pattern[i] == bm.Pattern[j]; i, j = i-1, j-1 {
		length++
	}
	return length
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
