package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"app/shared/database"
	"app/shared/server"

	"github.com/funayman/jmdict-api/install"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var (
	installFlag = flag.Bool("install", false, "install to a local db")
	configFlag  = flag.String("config", "config.json", "load a custom config file")
	config      = &configuration{}
)

type configuration struct {
	Database database.Setup `json:"database"`
	Install  install.Config `json:"installation"`
	Server   server.Server  `json:"server"`
}

func (c *configuration) Load(configPath string) {
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	//Get user input
	flag.Parse()

	//Load config
	log.Println("Loading config: " + *configFlag)
	config.Load(*configFlag)

	//Connect to database
	log.Println("Connecting to database...")
	database.Connect(config.Database)

	//Check if were installing database
	if *installFlag {
		err := install.JMDict(config.Install)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	//Do the other shit
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the Japanese Dictionary API\n")
	})

	r.HandleFunc("/word/{word}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		fmt.Fprintf(w, "Your word: %s", vars["word"])
	})

	server.Start(r, config.Server)

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
