package quransearch

import (
	"bufio"
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
	MIN_PATTERN_LEN    = 3
	MAX_INDEX_OF_LEN   = 10
	METHOD_INDEX_OF    = 0
	METHOD_BOYER_MOORE = 1
	METHOD_REGEX       = 2
	METHOD_DEFAULT     = METHOD_REGEX
)

var SurahName = [][]string{
	{"الفاتحة", "الفَاتِحَةِ"},
	{"البقرة", "البَقَرَةِ"},
	{"آل عمران", "آلِ عِمۡرَانَ"},
	{"النساء", "النِّسَاءِ"},
	{"المائدة", "المَائ‍ِدَةِ"},
	{"الأنعام", "الأَنعَامِ"},
	{"الأعراف", "الأَعۡرَافِ"},
	{"الأنفال", "الأَنفَالِ"},
	{"التوبة", "التَّوۡبَةِ"},
	{"يونس", "يُونُسَ"},
	{"هود", "هُودٍ"},
	{"يوسف", "يُوسُفَ"},
	{"الرعد", "الرَّعۡدِ"},
	{"إِبراهيم", "إِبۡرَاهِيمَ"},
	{"الحجر", "الحِجۡرِ"},
	{"النحل", "النَّحۡلِ"},
	{"الإِسراء", "الإِسۡرَاءِ"},
	{"الكهف", "الكَهۡفِ"},
	{"مريم", "مَرۡيَمَ"},
	{"طه", "طه"},
	{"الأنبياء", "الأَنبيَاءِ"},
	{"الحج", "الحَجِّ"},
	{"المؤمنون", "المُؤۡمِنُونَ"},
	{"النور", "النُّورِ"},
	{"الفرقان", "الفُرۡقَانِ"},
	{"الشعراء", "الشُّعَرَاءِ"},
	{"النمل", "النَّمۡلِ"},
	{"القصص", "القَصَصِ"},
	{"العنكبوت", "العَنكَبُوتِ"},
	{"الروم", "الرُّومِ"},
	{"لقمان", "لُقۡمَانَ"},
	{"السجدة", "السَّجۡدَةِ"},
	{"الأحزاب", "الأَحۡزَابِ"},
	{"سبأ", "سَبَإٍ"},
	{"فاطر", "فَاطِرٍ"},
	{"يس", "يسٓ"},
	{"الصافات", "الصَّافَّاتِ"},
	{"ص", "صٓ"},
	{"الزمر", "الزُّمَرِ"},
	{"غافر", "غَافِرٍ"},
	{"فصلت", "فُصِّلَتۡ"},
	{"الشورى", "الشُّورَىٰ"},
	{"الزخرف", "الزُّخۡرُفِ"},
	{"الدخان", "الدُّخَانِ"},
	{"الجاثية", "الجَاثِيةِ"},
	{"الأحقاف", "الأَحۡقَافِ"},
	{"محمد", "مُحَمَّدٍ"},
	{"الفتح", "الفَتۡحِ"},
	{"الحجرات", "الحُجُرَاتِ"},
	{"ق", "قٓ"},
	{"الذاريات", "الذَّارِيَاتِ"},
	{"الطور", "الطُّورِ"},
	{"النجم", "النَّجۡمِ"},
	{"القمر", "القَمَرِ"},
	{"الرحمن", "الرَّحۡمَٰن"},
	{"الواقعة", "الوَاقِعَةِ"},
	{"الحديد", "الحَدِيدِ"},
	{"المجادلة", "المُجَادلَةِ"},
	{"الحشر", "الحَشۡرِ"},
	{"الممتحنة", "المُمۡتَحنَةِ"},
	{"الصف", "الصَّفِّ"},
	{"الجمعة", "الجُمُعَةِ"},
	{"المنافقون", "المُنَافِقُونَ"},
	{"التغابن", "التَّغَابُنِ"},
	{"الطلاق", "الطَّلَاقِ"},
	{"التحريم", "التَّحۡرِيمِ"},
	{"الملك", "المُلۡكِ"},
	{"القلم", "القَلَمِ"},
	{"الحاقة", "الحَاقَّةِ"},
	{"المعارج", "المَعَارِجِ"},
	{"نوح", "نُوحٍ"},
	{"الجن", "الجِنِّ"},
	{"المزمل", "المُزَّمِّلِ"},
	{"المدثر", "المُدَّثِّرِ"},
	{"القيامة", "القِيَامَةِ"},
	{"الإِنسان", "الإِنسَانِ"},
	{"المرسلات", "المُرۡسَلَاتِ"},
	{"النبأ", "النَّبَإِ"},
	{"النازعات", "النَّازِعَاتِ"},
	{"عبس", "عَبَسَ"},
	{"التكوير", "التَّكۡوِيرِ"},
	{"الانفطار", "الانفِطَارِ"},
	{"المطففين", "المُطَفِّفِينَ"},
	{"الانشقاق", "الانشِقَاقِ"},
	{"البروج", "البُرُوجِ"},
	{"الطارق", "الطَّارِقِ"},
	{"الأعلى", "الأَعۡلَىٰ"},
	{"الغاشية", "الغَاشِيَةِ"},
	{"الفجر", "الفَجۡرِ"},
	{"البلد", "البَلَدِ"},
	{"الشمس", "الشَّمۡسِ"},
	{"الليل", "اللَّيۡلِ"},
	{"الضحى", "الضُّحَىٰ"},
	{"الشرح", "الشَّرۡحِ"},
	{"التين", "التِّينِ"},
	{"العلق", "العَلَقِ"},
	{"القدر", "القَدۡرِ"},
	{"البينة", "البَيِّنَةِ"},
	{"الزلزلة", "الزَّلۡزَلَةِ"},
	{"العاديات", "العَادِيَاتِ"},
	{"القارعة", "القَارِعَةِ"},
	{"التكاثر", "التَّكَاثُرِ"},
	{"العصر", "العَصۡرِ"},
	{"الهمزة", "الهُمَزَةِ"},
	{"الفيل", "الفِيلِ"},
	{"قريش", "قُرَيۡشٍ"},
	{"الماعون", "المَاعُونِ"},
	{"الكوثر", "الكَوثَرِ"},
	{"الكافرون", "الكَافِرُونَ"},
	{"النصر", "النَّصۡرِ"},
	{"المسد", "المَسَدِ"},
	{"الإِخلاص", "الإِخۡلَاصِ"},
	{"الفلق", "الفَلَقِ"},
	{"الناس", "النَّاسِ"},
}

type QuranSearch struct {
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

func (qs *QuranSearch) Search(p string, max int) []AyaMatch {
	if len(p) < MIN_PATTERN_LEN {
		return nil
	}

	if qs.oneLetterSpecialCase(p) || qs.twoLettersSpecialCase(p) {
		return qs.buildResults(qs.SpecialCases, len(p))
	}

	var matches []SearchMatch
	start := time.Now()

	switch qs.CurrentMethod {
	case METHOD_INDEX_OF:
		matches = qs.indexOfSearch(p, max, start)
	case METHOD_BOYER_MOORE:
		//matches = qs.boyerMooreSearch(p, max, start)
	case METHOD_REGEX:
		matches = qs.regexSearch(p, max, start)
	default:
		matches = qs.indexOfSearch(p, max, start)
	}

	return qs.buildResults(matches, len(p))
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
	var matches []SearchMatch
	re := regexp.MustCompile(p)
	for _, match := range re.FindAllStringIndex(qs.Quran, max) {
		matches = append(matches, *NewSearchMatch(qs.Quran, match[0], time.Since(start)))
	}
	return matches
}

func (qs *QuranSearch) buildResults(matches []SearchMatch, plen int) []AyaMatch {
	var results []AyaMatch
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
