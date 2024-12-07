package main

import (
	"fmt"
	"github.com/mbsoft31/quransearch/quransearch"
	"log"
)

func main() {
	qs, err := quransearch.NewQuranSearch("data/quran.txt")
	if err != nil {
		log.Fatal(err)
	}

	matches := qs.Search("الماكر", 10)
	for _, match := range matches {

		fmt.Printf("%v", match.BuildFullAya(qs.Quran))
	}

}
