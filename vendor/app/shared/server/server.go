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

func Start(r http.Handler, s Server) {
	logger.Info("webserver started on: " + s.address())
	http.ListenAndServe(s.address(), r)
}
