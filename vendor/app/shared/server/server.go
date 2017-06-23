package server

import (
	"fmt"
	"net/http"

	"app/shared/logger"
)

var (
	server *http.Server
)

type Server struct {
	Addr string `json:"host"`
	Port int    `json:"port"`
}

func (s Server) address() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

func Start(r http.Handler, s Server) {
	//set up our server
	server = &http.Server{Addr: s.address(), Handler: r}

	//start the webserver in its own go routine so it doesnt block main
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()

	//server started, log away
	logger.Info("webserver started on: " + s.address())
}

func Shutdown() {
	logger.Info("shutting down web server...")
	if err := server.Shutdown(nil); err != nil {
		logger.Fatal(err)
	}
}
