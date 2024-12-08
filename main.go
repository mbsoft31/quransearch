package main

import (
	"fmt"
	search "github.com/mbsoft31/quransearch/quransearch"
	"log"
)

func main() {
	qs, err := search.NewQuranSearch("data/quran.txt")
	if err != nil {
		log.Fatal(err)
	}

	matches := qs.Search("الناس", 100)

	var quranAr = &search.Quran{}
	var quranEn = &search.Quran{}

	err = search.ParseQuranXML("data/madina.xml", quranAr)
	if err != nil {
		log.Fatal(err)
	}
	err = search.ParseQuranXML("data/english-pickthall.xml", quranEn)
	if err != nil {
		log.Fatal(err)
	}
	for _, match := range matches {
		metaAr, _ := search.Fetch(quranAr, match)
		metaEn, _ := search.Fetch(quranEn, match)
		fmt.Printf("%s\n", match.StrBld.String())
		fmt.Printf("%v\n", metaAr.Aya.Text)
		fmt.Printf("%v\n", metaEn.Aya.Text)
	}

}
