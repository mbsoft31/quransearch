package quransearch

import (
	"encoding/xml"
	"strings"
	"time"
	"unicode/utf8"
)

// Structs and Models

// SearchMatch1 search match
type SearchMatch1 struct {
	Offset int
}

// SearchMethod is an interface for searching methods
type SearchMethod interface {
	Search(text, pattern string, max int) []SearchMatch
}

type indexOfMethod struct{}

func (i indexOfMethod) Search(text, pattern string, max int) []SearchMatch {
	start := time.Now()
	var matches []SearchMatch
	index := 0
	for len(matches) < max {
		println("check index = ", index)
		newIndex := strings.Index(text[index:], pattern)
		if newIndex == -1 {
			break
		}
		println("index = ", newIndex)
		matches = append(matches, *NewSearchMatch(text, newIndex, time.Since(start)))
		index += newIndex + utf8.RuneCountInString(pattern)
		println("new index = ", index)
	}
	return matches
}

// RegexMethod implements the SearchMethod interface
type RegexMethod struct{}

// BruteForceMethod implements the SearchMethod interface
type BruteForceMethod struct{}

// BoyerMooreMethod implements the Boyer-Moore string search algorithm
type BoyerMooreMethod struct {
	Pattern       string
	BadCharacter  []int
	GoodSuffix    []int
	PatternLength int
}

type SearchMatch struct {
	Index int
	Begin int
	End   int
	Word  int
	Surah int
	Aya   int
	Time  time.Duration
}

type AyaMatch struct {
	StrBld    strings.Builder
	Nfo       SearchMatch
	Len       int
	MLen      int
	SLen      int
	Indexes   []int
	PreSpaces int
}

type Quran struct {
	XMLName  xml.Name `xml:"quran"`
	Language string   `xml:"language,attr"`
	Version  string   `xml:"version,attr"`
	Source   string   `xml:"source,attr"`
	Surahs   []Surah  `xml:"surah"`
}

type Surah struct {
	No        int    `xml:"no,attr"`
	Name      string `xml:"name,attr"`
	Bismillah bool   `xml:"bismillah,attr"`
	Ayahs     []Ayah `xml:"ayat"`
}

type Ayah struct {
	No   int    `xml:"no,attr"`
	Text string `xml:"text,attr"`
}
