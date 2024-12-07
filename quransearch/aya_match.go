package quransearch

import (
	"fmt"
	"regexp"
	_ "regexp"
	"strconv"
	"strings"
	"time"
)

const uthmaniChars = "\u0650\u06e1\u0671\u0651\u064e\u0670\u064f\u0653\u06db\u0657\u0652\u06d6\u064c\u065e\u06e2\u06d7\u06e5\u0656\u06da\u06e6\u06de\u06d8\u064d\u200d\u0654\u064b\u06e7\u06dc\u06e0\u06e4\u06e9\u0655\u065c\u06ec\u06e8\u0640"
const dotsPrefix = "... "

// NewAyaMatch Constructor
func NewAyaMatch(quran string, ayaBegin bool, match SearchMatch, matchLen int) *AyaMatch {
	am := &AyaMatch{
		Nfo:     match,
		Indexes: []int{},
		MLen:    matchLen,
	}
	am.setNbrOfPreSpaces(quran)

	var oc int
	if !ayaBegin && match.Word > match.Begin {
		substr := quran[match.Word:match.End]
		am.StrBld.WriteString(dotsPrefix + substr)
		am.Len = len(dotsPrefix) + len(substr)
		oc = len(dotsPrefix) + (match.Index - match.Word)
	} else {
		substr := quran[match.Begin:match.End]
		am.StrBld.WriteString(substr)
		am.Len = len(substr)
		oc = match.Index - match.Begin
	}

	am.Indexes = append(am.Indexes, oc)
	return am
}

// Helper function to count preceding spaces
func (am *AyaMatch) setNbrOfPreSpaces(quran string) {
	if am.Nfo.Word > am.Nfo.Begin {
		am.PreSpaces = 1
		for i := am.Nfo.Begin + 1; i < am.Nfo.Word-1; i++ {
			if quran[i] == ' ' {
				am.PreSpaces++
			}
		}
	} else {
		am.PreSpaces = 0
	}
}

// AddOccurrence Add a new occurrence index
func (am *AyaMatch) AddOccurrence(next int) {
	am.Indexes = append(am.Indexes, next)
}

// AppendNumber Append a suffix to the aya text
func (am *AyaMatch) AppendNumber(suffix string) {
	am.StrBld.WriteString(suffix)
	am.SLen = len(suffix)
	am.Len += am.SLen
}

// BuildUthmaniRegEx Build a Uthmani regex pattern
func (am *AyaMatch) BuildUthmaniRegEx() string {
	var b strings.Builder
	p := am.StrBld.String()[am.Indexes[0] : am.Indexes[0]+am.MLen]

	for i, c := range p {
		switch c {
		case 'آ':
			b.WriteString("([أا]|ءا|ءؤ)?")
		case 'ا':
			b.WriteString("[اؤوى]?")
		case 'ي':
			b.WriteString("ا?[يأء]?")
		case 'ئ':
			b.WriteString("[ئءإ]?ي?")
		case 'أ':
			b.WriteString("([أئءاي]?|ؤا)")
		case 'ء':
			b.WriteString("[ءيا]?")
		case 'إ':
			b.WriteString("[إء]ي?")
		case 'ؤ':
			b.WriteString("[ؤء]")
		case 'ى':
			b.WriteString("[ىا]")
		case 'و':
			b.WriteString("[وا]?ا?")
		case 'س':
			b.WriteString("[سص]")
		case 'ش':
			b.WriteString("شا?")
		case 'ل':
			if i > 1 && i+1 < len(p) &&
				rune(p[i-1]) == 'ل' && rune(p[i-2]) == 'ا' &&
				(rune(p[i+1]) == 'ي' || rune(p[i+1]) == 'ا' || rune(p[i+1]) == 'ذ') {
				b.WriteString("ل?")
			} else {
				b.WriteRune(c)
			}
		case 'ن':
			if (i > 1 && rune(p[i-1]) == 'أ' && rune(p[i-2]) == 'و') || (i > 0 && rune(p[i-1]) == 'ن') {
				b.WriteString("ن?")
			} else {
				b.WriteRune(c)
			}
		case 'ة':
			b.WriteString("[ةت]")
		case ' ':
			if i > 1 && rune(p[i-1]) == 'ا' && (rune(p[i-2]) == 'ي' || rune(p[i-2]) == 'ه' || rune(p[i-2]) == 'م') {
				b.WriteString(" ?")
			} else {
				b.WriteRune(c)
			}
			continue
		default:
			b.WriteRune(c)
		}
		b.WriteRune('[')
		b.WriteString(uthmaniChars)
		b.WriteString("]{0,5}")
	}
	return b.String()
}

func (am *AyaMatch) BuildFullAya(quran string) string {
	var b strings.Builder
	b.WriteString(quran[am.Nfo.Begin:am.Nfo.End])
	return b.String()
}

// NewSearchMatch Constructor for SearchMatch
func NewSearchMatch(quran string, i int, t time.Duration) *SearchMatch {
	sm := &SearchMatch{
		Index: i,
		Time:  t,
	}

	// Parse input for the other info
	sm.Begin = strings.LastIndex(quran[:i], "\n") + 1 // covers -1 too

	indexNext := sm.setSurahNumber(quran)

	sm.setAyaNumber(quran, indexNext)

	sm.Word = strings.LastIndex(quran[:i], " ") + 1
	if sm.Word == 0 || sm.Word < sm.Begin {
		sm.Word = sm.Begin
	}

	sm.End = strings.Index(quran[i:], "\n") + i
	// quran.txt does end with \n before EOF

	return sm
}

func (sm *SearchMatch) setAyaNumber(quran string, indexNext int) {
	sm.Begin = strings.Index(quran[indexNext:], "|") + indexNext
	sm.Aya, _ = strconv.Atoi(quran[indexNext:sm.Begin])
	sm.Begin++
}

func (sm *SearchMatch) setSurahNumber(quran string) int {
	n := strings.Index(quran[sm.Begin:], "|") + sm.Begin
	sm.Surah, _ = strconv.Atoi(quran[sm.Begin:n])
	n++
	return n
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
	fmt.Printf("in %4d micro-sec, at %s\n", sm.Time.Microseconds(), res)
}
