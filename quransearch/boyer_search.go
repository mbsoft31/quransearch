package quransearch

import (
	"sync"
	"unsafe"
)

// BoyerMoore represents the optimized Boyer-Moore search algorithm.
type BoyerMoore struct {
	pattern         []byte
	badCharTable    []int
	goodSuffixTable []int
	patternLen      int
}

// NewBoyerMoore initializes a BoyerMoore instance.
func NewBoyerMoore(pattern string) *BoyerMoore {
	p := []byte(pattern)
	return &BoyerMoore{
		pattern:         p,
		badCharTable:    makeBadCharTable(p),
		goodSuffixTable: makeGoodSuffixTable(p),
		patternLen:      len(p),
	}
}

// Search performs the Boyer-Moore search.
func (bm *BoyerMoore) Search(text string, max int) []SearchMatch1 {
	if max <= 0 {
		max = int(^uint(0) >> 1) // Set max to the largest int value.
	}

	textBytes := *(*[]byte)(unsafe.Pointer(&text)) // Direct byte operations.
	textLen := len(textBytes)
	if bm.patternLen == 0 || textLen < bm.patternLen {
		return nil
	}

	matches := make([]SearchMatch1, 0, max)
	lastIndex := bm.patternLen - 1
	i := lastIndex

	for i < textLen && len(matches) < max {
		j := lastIndex
		for j >= 0 && bm.pattern[j] == textBytes[i] {
			i--
			j--
		}
		if j < 0 { // Full match found.
			matches = append(matches, SearchMatch1{Offset: i + 1})
			i += bm.patternLen * 2 // Skip overlapping matches.
		} else {
			i += int(maxInt(bm.goodSuffixTable[lastIndex-j], bm.badCharTable[textBytes[i]]))
		}
	}
	return matches
}

// makeBadCharTable constructs the bad character shift table.
func makeBadCharTable(pattern []byte) []int {
	alphabetSize := 256 // Extended ASCII.
	table := make([]int, alphabetSize)
	patternLen := len(pattern)

	for i := 0; i < alphabetSize; i++ {
		table[i] = patternLen
	}
	for i := 0; i < patternLen-1; i++ {
		table[pattern[i]] = patternLen - 1 - i
	}
	return table
}

// makeGoodSuffixTable constructs the good suffix shift table.
func makeGoodSuffixTable(pattern []byte) []int {
	patternLen := len(pattern)
	table := make([]int, patternLen)
	lastPrefixPos := patternLen

	for i := patternLen - 1; i >= 0; i-- {
		if isPrefix(pattern, i+1) {
			lastPrefixPos = i + 1
		}
		table[patternLen-1-i] = lastPrefixPos - i + patternLen - 1
	}
	for i := 0; i < patternLen-1; i++ {
		slen := suffixLength(pattern, i)
		table[slen] = patternLen - 1 - i + slen
	}
	return table
}

// isPrefix checks if pattern[p:] is a prefix of pattern.
func isPrefix(pattern []byte, p int) bool {
	patternLen := len(pattern)
	for i, j := p, 0; i < patternLen; i, j = i+1, j+1 {
		if pattern[i] != pattern[j] {
			return false
		}
	}
	return true
}

// suffixLength calculates the maximum length of a substring ending at p that is also a suffix.
func suffixLength(pattern []byte, p int) int {
	length := 0
	for i, j := p, len(pattern)-1; i >= 0 && pattern[i] == pattern[j]; i, j = i-1, j-1 {
		length++
	}
	return length
}

// maxInt returns the larger of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ParallelSearch performs parallel Boyer-Moore search for large texts.
func (bm *BoyerMoore) ParallelSearch(text string, max int) []SearchMatch1 {
	const chunkSize = 10000
	textBytes := *(*[]byte)(unsafe.Pointer(&text))
	textLen := len(textBytes)
	numChunks := (textLen + chunkSize - 1) / chunkSize

	var wg sync.WaitGroup
	matchesChan := make(chan []SearchMatch1, numChunks)

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > textLen {
			end = textLen
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			matchesChan <- bm.Search(string(textBytes[start:end]), max)
		}(start, end)
	}

	wg.Wait()
	close(matchesChan)

	var matches []SearchMatch1
	for matchList := range matchesChan {
		matches = append(matches, matchList...)
		if len(matches) >= max {
			break
		}
	}
	return matches
}
