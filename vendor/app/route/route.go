package route

import (
	"app/shared/router"
	"net/http"
)

func Load() http.Handler {
	return middleware(router.Instance())
}

func middleware(h http.Handler) http.Handler {
	return h
}
