package main

import (
	"encoding/xml"
	"fmt"
	"github.com/mbsoft31/quransearch/quransearch"
	"log"
	"os"
)

func main() {
	qs, err := quransearch.NewQuranSearch("data/quran.txt")
	if err != nil {
		log.Fatal(err)
	}

	matches := qs.Search("الماكر", 10)

	quranAr, err := parseQuranXML("data/madina.xml")
	quranEn, err := parseQuranXML("data/english-pickthall.xml")

	for _, match := range matches {
		fmt.Printf("%v\n", match.Nfo.Time)
		fmt.Printf("%v ", qs.GetAyaPrefix(match.Nfo.Surah, match.Nfo.Aya))
		fmt.Printf("%v\n", qs.Quran[match.Nfo.Begin:match.Nfo.End])
		//fmt.Printf("%v\n", qs.GetAyaSuffix(match.Nfo.Surah, match.Nfo.Ayat))
		metaAr, _ := fetch(quranAr, match)
		metaEn, _ := fetch(quranEn, match)
		fmt.Printf("%v\n", metaAr.Ayat.Text)
		fmt.Printf("%v\n", metaEn.Ayat.Text)
	}

}

type MetaData struct {
	Surah Surah
	Ayat  Ayat
}

func fetch(quran *Quran, m quransearch.AyaMatch) (*MetaData, error) {
	sura := quran.Surahs[m.Nfo.Surah-1]
	return &MetaData{
		Surah: quran.Surahs[m.Nfo.Surah-1],
		Ayat:  sura.Ayats[m.Nfo.Aya-1],
	}, nil
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
	Ayats     []Ayat `xml:"ayat"`
}

type Ayat struct {
	No   int    `xml:"no,attr"`
	Text string `xml:"text,attr"`
}

func parseQuranXML(filename string) (*Quran, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var quran Quran
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&quran)
	if err != nil {
		return nil, err
	}

	return &quran, nil
}
