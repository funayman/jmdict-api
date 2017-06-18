package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

var (
	r RouteHolder
)

type RouteHolder struct {
	*mux.Router
}

func init() {
	r.Router = mux.NewRouter()
}

// ReadConfig returns the information
func ReadConfig() RouteHolder {
	return r
}

// Instance returns the router
func Instance() *mux.Router {
	return r.Router
}

func Route(path string, fn http.HandlerFunc) {
	r.HandleFunc(path, fn)
}

func GetParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}
