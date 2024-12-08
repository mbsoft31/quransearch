package quransearch

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type MetaData struct {
	Surah Surah
	Aya   Ayah
}

func Fetch(quran *Quran, m AyaMatch) (*MetaData, error) {
	sura := quran.Surahs[m.Nfo.Surah-1]
	return &MetaData{
		Surah: quran.Surahs[m.Nfo.Surah-1],
		Aya:   sura.Ayahs[m.Nfo.Aya-1],
	}, nil
}

func ParseQuranXML(filename string, quran *Quran) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ParseQuranXML: %v", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal("ParseQuranXML: Could not close file!")
		}
	}(file)

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&quran)
	if err != nil {
		return fmt.Errorf("ParseQuranXML: %v", err)
	}

	return nil
}
