package quransearch

import (
	"strings"
	"time"
)

// Structs and Models

// SearchMatch1 search match
type SearchMatch1 struct {
	Offset int
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
