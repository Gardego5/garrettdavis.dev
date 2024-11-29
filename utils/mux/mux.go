package mux

import (
	"net/http"
	"strings"
)

type ServeMux struct {
	// middleware stack
	stack
	// inner is the inner http.ServeMux that this ServeMux wraps.
	inner  *http.ServeMux
	prefix string
}

func NewServeMux(
	register func(m *ServeMux),
	middleware ...Middleware,
) *ServeMux {
	mux := &ServeMux{
		stack: *new(stack).with(middleware...),
		inner: http.NewServeMux(), prefix: "",
	}
	if register != nil {
		register(mux)
	}
	return mux
}

func (mux *ServeMux) Handle(
	pattern string,
	handler http.Handler,
	middleware ...Middleware,
) {
	var method, path string
	prefix := mux.prefix

	// Split the pattern into a method and path if it exists.
	parts := strings.SplitN(pattern, " ", 2)
	switch len(parts) {
	case 1:
		switch parts[0] {
		case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
			method = parts[0]
		default:
			path = parts[0]
		}
	case 2:
		method, path = parts[0], parts[1]
	}

	// If the prefix ends with a "/" and the path starts with a "/" then we need
	// to remove one of them to avoid a double slash.
	if strings.HasSuffix(prefix, "/") && strings.HasPrefix(path, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	fullPath := prefix + path

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

	stk := &mux.stack
	if len(middleware) > 0 {
		stk = stk.with(middleware...)
	}
	handler = stk.Use(handler)
	mux.inner.Handle(method+fullPath, handler)
}

func (mux *ServeMux) HandleFunc(
	pattern string,
	handler func(http.ResponseWriter, *http.Request),
	middleware ...Middleware,
) {
	mux.Handle(pattern, http.HandlerFunc(handler), middleware...)
}

func (mux *ServeMux) Use(
	register func(m *ServeMux),
	middleware ...Middleware,
) *ServeMux {
	child := &ServeMux{
		stack: *mux.stack.with(middleware...),
		inner: mux.inner, prefix: mux.prefix,
	}
	if register != nil {
		register(child)
	}
	return child
}

func (mux *ServeMux) Group(
	prefix string,
	register func(m *ServeMux),
	middleware ...Middleware,
) *ServeMux {
	child := &ServeMux{
		stack: *mux.stack.with(middleware...),
		inner: mux.inner, prefix: mux.prefix + prefix,
	}
	if register != nil {
		register(child)
	}
	return child
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.inner.ServeHTTP(w, r)
}
