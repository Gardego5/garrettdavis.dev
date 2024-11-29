package mux

import "net/http"

type (
	// this is a middleware always present in the top level of the ServeMux.
	// it acts as a container for all other middleware, and data that is passe
	// to each request
	stack struct {
		middleware []Middleware
		parent     *stack
	}

	MiddlewareFunc func(next http.Handler) http.Handler
	Middleware     interface {
		Use(next http.Handler) http.Handler
	}
)

// assertions to confirm that the interfaces are implemented
var _ Middleware = (*stack)(nil)
var _ Middleware = MiddlewareFunc(nil)

func (f MiddlewareFunc) Use(next http.Handler) http.Handler {
	return f(next)
}

func (stk *stack) Use(next http.Handler) http.Handler {
	// compose the middleware in reverse order so that middlewares added later
	// are able to access the context data added by middlewares added earlier
	for i := len(stk.middleware) - 1; i >= 0; i-- {
		next = stk.middleware[i].Use(next)
	}

	return next
}

// creates a new stack struct with the given middleware & datasources
func (stk *stack) with(middleware ...Middleware) (child *stack) {
	child = new(stack)
	child.parent = stk

	// copy the existing middleware
	child.middleware = make([]Middleware, 0, len(stk.middleware))
	child.middleware = append(child.middleware, stk.middleware...)

	// add the new middleware
	for _, middle := range middleware {
		if middle, ok := middle.(Middleware); ok {
			child.middleware = append(child.middleware, middle)
		}
	}

	return child
}
