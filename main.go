package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/funayman/jmdict-api/install"
	_ "github.com/mattn/go-sqlite3"
)

var (
	installFlag = flag.Bool("install", false, "install to a local db")
	configFlag  = flag.String("config", "./config.json", "load a custom config file")
)

func main() {
	flag.Parse()
	if *installFlag {
		err := install.JMDict()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	fmt.Println("This is the app T_T")
	/*
		//kanjidic2
		kanjiFile, err := os.Open("./data/kanjidic2.xml")
		if err != nil {
			log.Fatal(err)
		}
		defer kanjiFile.Close()

		kanji, err := install.LoadKanjiDic2(kanjiFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Len of Kanji: %d\n", len(kanji))
		fmt.Printf("%+v\n", kanji[100])
	*/
}
