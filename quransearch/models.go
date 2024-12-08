package quransearch

import (
	"encoding/xml"
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
