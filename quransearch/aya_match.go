package quransearch

import (
	"fmt"
	_ "regexp"
	"strings"
)

const uthmaniChars = "\u0650\u06e1\u0671\u0651\u064e\u0670\u064f\u0653\u06db\u0657\u0652\u06d6\u064c\u065e\u06e2\u06d7\u06e5\u0656\u06da\u06e6\u06de\u06d8\u064d\u200d\u0654\u064b\u06e7\u06dc\u06e0\u06e4\u06e9\u0655\u065c\u06ec\u06e8\u0640"
const dotsPrefix = "... "

type AyaMatch struct {
	StrBld    strings.Builder
	Nfo       SearchMatch
	Len       int
	MLen      int
	SLen      int
	Indexes   []int
	PreSpaces int
}

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

// Build a Uthmani regex pattern
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
	b.WriteString(fmt.Sprintf(" {%s %d}\n", SurahName[am.Nfo.Surah][0], am.Nfo.Aya))
	return b.String()
}
