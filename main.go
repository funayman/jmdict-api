package main

import (
	"encoding/json"
	"flag"
	"os"
	"runtime"

	"app/controller"
	"app/install"
	"app/route"
	"app/shared/database"
	"app/shared/logger"
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
	Log      logger.Config  `json:"logger"`
}

func (c *configuration) Load(configPath string) {
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	err = json.NewDecoder(configFile).Decode(c)
	if err != nil {
		panic(err)
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
	config.Load(*configFlag)

	//Start my logger
	logger.Load(config.Log)

	//Connect to database
	logger.Info("Connecting to database...")
	database.Connect(config.Database)

	//Check if were installing database
	if *installFlag {
		logger.Info("Installing JMDict...")
		err := install.JMDict(config.Install)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("Installing KanjiDic2...")
		err = install.KanjiDic2(config.Install)
		if err != nil {
			logger.Fatal(err)
		}

		os.Exit(0)
	}

	//Load the controller routes
	logger.Info("Loading controllers...")
	controller.Load()

	//Load all routes and middleware and start the server
	logger.Info("Starting web server...")
	server.Start(route.Load(), config.Server)

	//TODO cleanup on signal intterup
	/*
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				logger.Info("Received an interrupt, stopping services...")
				logger.Close()
				cleanupDone <- true
			}
		}()
		<-cleanupDone
	*/
}
