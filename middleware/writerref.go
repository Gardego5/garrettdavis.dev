package middleware

import (
	"context"
	"net/http"
)

// This middleware adds the http.ResponseWriter to the request context.
// writerrefKey is the key used to store the http.ResponseWriter in the request context.
// This allows for utility functions that write to the request without explicitly passing the http.ResponseWriter.
//
// Example:
//
//	func GetPage(w http.ResponseWriter, r *http.Request) {
//		middleware.RenderPage(r, Title{"Hello, World!"}, Div{"Hello, World!"})
//	}
func WriterRef(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, writerrefKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
