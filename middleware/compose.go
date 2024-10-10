package middleware

import (
	"net/http"

	"github.com/elliotchance/pie/v2"
)

type Middleware func(http.Handler) http.Handler

func Compose(middleware ...Middleware) Middleware {
	return pie.Reduce(middleware, func(one, two Middleware) Middleware {
		return func(handler http.Handler) http.Handler {
			return one(two(handler))
		}
	})
}
