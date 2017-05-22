package main

import (
	"fmt"
	"log"
	"os"

	"github.com/funayman/jmdict-api/parser"
)

func main() {
	//get the file
	data, err := os.Open("./data/JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	words, err := parser.LoadJMDict(data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Current count: ", len(words))

	//kanjidic2
	kanjiFile, err := os.Open("./data/kanjidic2.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	kanji, err := parser.LoadKanjiDic2(kanjiFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Len of Kanji: %d\n", len(kanji))
	fmt.Printf("%+v\n", kanji[100])
}
