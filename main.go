package main

import (
	"embed"
	"fmt"
	search "github.com/mbsoft31/quransearch/quransearch"
	"log"
)

//go:embed data/quran.txt
var quranText embed.FS

func main() {
	qs, err := search.NewQuranSearch("data/quran.txt")
	// qs, err := search.NewQuranSearchWithText(quranText)
	if err != nil {
		log.Fatal(err)
	}

	matches := qs.Search("ناس", 10)

	/*var quranAr = &search.Quran{}
	var quranEn = &search.Quran{}

	err = search.ParseQuranXML("data/madina.xml", quranAr)
	if err != nil {
		log.Fatal(err)
	}
	err = search.ParseQuranXML("data/english-pickthall.xml", quranEn)
	if err != nil {
		log.Fatal(err)
	}*/
	for index, match := range matches {
		//metaAr, _ := search.Fetch(quranAr, match)
		//metaEn, _ := search.Fetch(quranEn, match)
		//fmt.Printf("%s\n", match.StrBld.String())
		//fmt.Printf("%v\n", metaAr.Aya.Text)
		//fmt.Printf("%v\n", metaEn.Aya.Text)
		fmt.Print("index: ", index, " => ")
		match.Nfo.Print()
		fmt.Printf("%s\n", match.StrBld.String())
	}

}
