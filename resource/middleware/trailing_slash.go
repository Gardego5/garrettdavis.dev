package middleware

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/utils/mux"
)

var TrailingSlash mux.Middleware = mux.MiddlewareFunc(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
})
