package quransearch

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	DEF_SEARCH_LIMIT   = 10
	DEF_BUFFER_SIZE    = 1024 * 4
	MIN_PATTERN_LEN    = 1
	MAX_INDEX_OF_LEN   = 10
	METHOD_INDEX_OF    = 0
	METHOD_BOYER_MOORE = 1
	METHOD_REGEX       = 2
	METHOD_BRUTE_FORCE = 3
	METHOD_DEFAULT     = METHOD_REGEX
)

/*
	TODO: Fix the Boyer Moore and IndexOf methods they giving wrong results
*/

type QuranSearch struct {
	Reader        *bufio.Reader
	Quran         string
	CurrentMethod int
	SurahAyaNbrs  bool
	AyaBegin      bool
	SpecialCases  []SearchMatch
}

func NewQuranSearch(filePath string) (*QuranSearch, error) {
	qs := &QuranSearch{CurrentMethod: METHOD_DEFAULT}
	err := qs.readFile(filePath)
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func NewQuranSearchWithText(quranFile embed.FS) (*QuranSearch, error) {
	qs := &QuranSearch{CurrentMethod: METHOD_DEFAULT}
	file, err := quranFile.ReadFile("data/quran.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading quran file: %v", err)
	}
	qs.Quran = string(file)
	return qs, nil
}

func (qs *QuranSearch) readFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Could not close file: %s", err.Error())
		}
	}(file)

	qs.Reader = bufio.NewReader(file)
	// TOOD: use the reader instead of readind all file line by line

	var sb strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sb.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	qs.Quran = sb.String()
	return nil
}

func (qs *QuranSearch) Search(p string, max int) []AyaMatch {
	if len(p) < MIN_PATTERN_LEN {
		return nil
	}

	if qs.oneLetterSpecialCase(p) || qs.twoLettersSpecialCase(p) {
		return qs.buildResults(qs.SpecialCases, len(p))
	}

	var matches []SearchMatch

	switch qs.CurrentMethod {
	case METHOD_BOYER_MOORE:
		boyerMoore := BoyerMooreMethod{}
		matches = boyerMoore.Search(qs.Quran, p, max)
	case METHOD_REGEX:
		regex := RegexMethod{}
		matches = regex.Search(qs.Quran, p, max)
	case METHOD_BRUTE_FORCE:
		bruteForce := BruteForceMethod{}
		matches = bruteForce.Search(qs.Quran, p, max)
	case METHOD_INDEX_OF:
		indexOf := indexOfMethod{}
		matches = indexOf.Search(qs.Quran, p, max)
	default:
		indexOf := indexOfMethod{}
		matches = indexOf.Search(qs.Quran, p, max)
	}

	return qs.buildResults(matches, len(p))
}

func (qs *QuranSearch) oneLetterSpecialCase(p string) bool {
	var pchar = int32(p[0])
	switch pchar {
	case 'ص':
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 335061, 0)}
	case 'ق':
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 384642, 0)}
	case 'ن':
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 421495, 0)}
	default:
		return false
	}
	return true
}

func (qs *QuranSearch) twoLettersSpecialCase(p string) bool {
	switch p {
	case "طه":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 227524, 0)}
	case "طس":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 277440, 0)}
	case "يس":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 324531, 0)}
	case "ص ":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 335061, 0)}
	case "حم":
		qs.SpecialCases = []SearchMatch{
			*NewSearchMatch(qs.Quran, 346076, 0),
			*NewSearchMatch(qs.Quran, 353019, 0),
			*NewSearchMatch(qs.Quran, 357570, 0),
			*NewSearchMatch(qs.Quran, 362337, 0),
			*NewSearchMatch(qs.Quran, 367420, 0),
			*NewSearchMatch(qs.Quran, 369667, 0),
			*NewSearchMatch(qs.Quran, 372513, 0),
		}
	case "ق ":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 384642, 0)}
	case "ن ":
		qs.SpecialCases = []SearchMatch{*NewSearchMatch(qs.Quran, 421495, 0)}
	default:
		return false
	}
	return true
}

func (qs *QuranSearch) indexOfSearch(p string, max int, start time.Time) []SearchMatch {
	var matches []SearchMatch
	index := 0
	for len(matches) < max {
		index = strings.Index(qs.Quran[index:], p)
		if index == -1 {
			break
		}
		matches = append(matches, *NewSearchMatch(qs.Quran, index, time.Since(start)))
		index += len(p)
	}
	return matches
}

func (qs *QuranSearch) regexSearch(p string, max int, start time.Time) []SearchMatch {
	var matches = make([]SearchMatch, 0)
	re := regexp.MustCompile(p)
	for _, match := range re.FindAllStringIndex(qs.Quran, max) {
		matches = append(matches, *NewSearchMatch(qs.Quran, match[0], time.Since(start)))
	}
	return matches
}

func (qs *QuranSearch) buildResults(matches []SearchMatch, plen int) []AyaMatch {
	var results = make([]AyaMatch, 0)
	for _, match := range matches {
		results = append(results, *NewAyaMatch(qs.Quran, qs.AyaBegin, match, plen))
	}
	return results
}

func (qs *QuranSearch) GetAyaSuffix(surah, aya int) string {
	return fmt.Sprintf(" \u200F[%s %d]", SurahName[surah-1][0], aya)
}

func (qs *QuranSearch) GetAyaPrefix(surah, aya int) string {
	return fmt.Sprintf("[%s %d] ", SurahName[surah-1][0], aya)
}
