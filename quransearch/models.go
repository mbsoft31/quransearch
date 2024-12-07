package quransearch

import (
	"fmt"
	"regexp"
	"strconv"
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

// NewSearchMatch Constructor for SearchMatch
func NewSearchMatch(quran string, i int, t time.Duration) *SearchMatch {
	sm := &SearchMatch{
		Index: i,
		Time:  t,
	}

	// Parse input for the other info
	sm.Begin = strings.LastIndex(quran[:i], "\n") + 1 // covers -1 too

	n := strings.Index(quran[sm.Begin:], "|") + sm.Begin
	sm.Surah, _ = strconv.Atoi(quran[sm.Begin:n])
	n++

	sm.Begin = strings.Index(quran[n:], "|") + n
	sm.Aya, _ = strconv.Atoi(quran[n:sm.Begin])
	sm.Begin++

	sm.Word = strings.LastIndex(quran[:i], " ") + 1
	if sm.Word == 0 || sm.Word < sm.Begin {
		sm.Word = sm.Begin
	}

	sm.End = strings.Index(quran[i:], "\n") + i
	// quran.txt does end with \n before EOF

	return sm
}

// NewSearchMatchWithoutTime Overloaded constructor without time
func NewSearchMatchWithoutTime(quran string, i int) *SearchMatch {
	return NewSearchMatch(quran, i, 0)
}

// NewSearchMatchFromMatcher Constructor using regex Matcher
func NewSearchMatchFromMatcher(m *regexp.Regexp, input string, t time.Duration) *SearchMatch {
	sm := &SearchMatch{
		Time: t,
	}

	// RegEx matching patterns are always limited to '\|([0-9]+)\|([0-9]+)\|(.*(pattern).*)\n'
	matches := m.FindStringSubmatchIndex(input)
	if matches == nil {
		return nil
	}

	sm.Surah, _ = strconv.Atoi(input[matches[2]:matches[3]])
	sm.Aya, _ = strconv.Atoi(input[matches[4]:matches[5]])

	sm.Begin = matches[6]
	sm.End = matches[7] + 1

	sm.Index = matches[8]

	sm.Word = strings.LastIndex(input[:sm.Index], " ") + 1
	if sm.Word == 0 || sm.Word < sm.Begin {
		sm.Word = sm.Begin
	}

	return sm
}

// NewSearchMatchFromMatcherWithoutTime Overloaded constructor without time using regex Matcher
func NewSearchMatchFromMatcherWithoutTime(m *regexp.Regexp, input string) *SearchMatch {
	return NewSearchMatchFromMatcher(m, input, 0)
}

// NewSearchMatchCopy Copy constructor
func NewSearchMatchCopy(s *SearchMatch) *SearchMatch {
	return &SearchMatch{
		Index: s.Index,
		Time:  s.Time,
		Begin: s.Begin,
		End:   s.End,
		Word:  s.Word,
		Surah: s.Surah,
		Aya:   s.Aya,
	}
}

// Print method
func (sm *SearchMatch) Print() {
	res := fmt.Sprintf("[%d،%d] quran[b=%d; w=%d; i=%d; e=%d].", sm.Surah, sm.Aya, sm.Begin, sm.Word, sm.Index, sm.End)
	fmt.Printf("in %8d micro-sec, at %s\n", sm.Time.Microseconds(), res)
}
