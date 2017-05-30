package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

func (s Server) address() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

func Start(r http.Handler, s Server) {
	log.Println("starting webserver on:" + s.address())
	http.ListenAndServe(s.address(), r)
}
