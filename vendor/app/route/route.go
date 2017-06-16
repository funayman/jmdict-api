package route

import (
	"net/http"

	"app/route/middleware/logrequest"
	"app/shared/router"
)

func Load() http.Handler {
	return middleware(router.Instance())
}

func middleware(h http.Handler) http.Handler {
	h = logrequest.Handler(h)

	return h
}
