package middleware

import (
	"net/http"
	"strings"

	"github.com/Gardego5/garrettdavis.dev/utils/mux"
)

var IgnoreIndex = mux.MiddlewareFunc(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
})
