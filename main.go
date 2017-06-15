package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"runtime"

	"app/controller"
	"app/install"
	"app/route"
	"app/shared/database"
	"app/shared/server"

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

	// Load the controller routes
	controller.Load()

	server.Start(route.Load(), config.Server)
}
