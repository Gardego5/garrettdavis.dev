package middleware

import (
	"net/http"
	"strings"
)

type ServeMux struct {
	inner      *http.ServeMux
	middleware []Middleware
	prefix     string
}

func NewServeMux(middleware ...Middleware) *ServeMux {
	return &ServeMux{
		inner:      http.NewServeMux(),
		middleware: middleware,
		prefix:     "",
	}
}

func (mux *ServeMux) compose(handler http.Handler) http.Handler {
	middleware := Compose(mux.middleware...)
	return middleware(handler)
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler) {
	var method, path string
	prefix := mux.prefix

	parts := strings.SplitN(pattern, " ", 2)
	switch len(parts) {
	case 1:
		path = parts[0]
	case 2:
		method, path = parts[0], parts[1]
	}

	// Add a space after the method if it exists.
	// This will result in a valid pattern when reassembled:
	// ```
	//	prefix := "/foo"
	//	pattern := "GET /bar"
	//	method, path := "GET", "/bar"
	//	method = "GET " // <- add space
	//	// Reassemble
	//	pattern = method + prefix + path
	//	pattern = "GET /foo/bar"
	// ```
	if len(method) > 0 {
		method = method + " "
	}

	// If the prefix ends with a "/" and the path starts with a "/" then we
	// need to remove one of them to avoid a double slash.
	if strings.HasSuffix(prefix, "/") && strings.HasPrefix(path, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	mux.inner.Handle(method+prefix+path, mux.compose(handler))
}

func (mux *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(pattern, http.HandlerFunc(handler))
}

func (mux *ServeMux) Use(middleware ...Middleware) *ServeMux {
	m := make([]Middleware, len(mux.middleware))
	copy(m, mux.middleware)

	return &ServeMux{
		inner:      mux.inner,
		middleware: append(m, middleware...),
		prefix:     mux.prefix,
	}
}

func (mux *ServeMux) Group(prefix string, middleware ...Middleware) *ServeMux {
	m := make([]Middleware, len(mux.middleware))
	copy(m, mux.middleware)

	return &ServeMux{
		inner:      mux.inner,
		middleware: append(m, middleware...),
		prefix:     mux.prefix + prefix,
	}
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.inner.ServeHTTP(w, r)
}
