package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"app/install"
	"app/model"
	"app/shared/database"
	"app/shared/router"
	"app/shared/server"

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
		log.Println("Installing JMDict...")
		err := install.JMDict(config.Install)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Installing KanjiDic2...")
		err = install.KanjiDic2(config.Install)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	//Do the other shit
	router.Route("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to the Japanese Dictionary API\n")
	})

	router.Route("/word/{word}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		word := vars["word"]
		queries := r.URL.Query()
		words := []*model.Word{}

		//grab the base ID for the word(s)
		rows, err := database.SQL.Query(database.QuerySearchForID, word, word, word)
		defer rows.Close()

		if err != nil {
			log.Print(err)
		}

		var id int
		for rows.Next() {
			rows.Scan(&id)
			words = append(words, &model.Word{ID: id})
		}

		for _, word := range words {
			err = word.BuildSelf()
			if err != nil {
				log.Println(err)
			}
		}

		//write the output
		var output []byte
		switch f := queries.Get("format"); strings.ToLower(f) {
		case "xml":
			output, err = xml.Marshal(words)
			if err != nil {
				log.Println(err)
			}
		default:
			output, err = json.Marshal(words)
			if err != nil {
				log.Println(err)
			}
		}

		fmt.Fprintf(w, "%s\n", output)
	})

	server.Start(config.Server)
}
