package server

import (
	"fmt"
	"net/http"

	"app/shared/logger"
)

type Server struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func (s Server) address() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

//TODO graceful startup and shutdown (go 1.8 server.Shutdown())

func Start(r http.Handler, s Server) {
	logger.Info("webserver started on: " + s.address())
	logger.Fatal(http.ListenAndServe(s.address(), r))
}
