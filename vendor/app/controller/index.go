package controller

import (
	"app/shared/router"
	"fmt"
	"net/http"
)

func init() {
	router.Route("/", IndexFunc)
}

func IndexFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Japanese Dictionary API\n")
}
