package controller

import (
	"net/http"

	"app/shared/router"
)

func init() {
	router.Route("/", IndexFunc)
}

func IndexFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Japanese Dictionary API\n"))
}
