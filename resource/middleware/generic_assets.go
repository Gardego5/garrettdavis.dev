package middleware

import (
	"context"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/resource/internal"
	"github.com/Gardego5/garrettdavis.dev/utils/mux"
)

type (
	boundKey[A any] struct {
		key internal.GenericKey[A]
		val *A
	}
	syringe interface {
		inject(context.Context) context.Context
	}
)

func (b boundKey[A]) inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, b.key, b.val)
}

func Syringe[A any](a *A) boundKey[A] {
	return boundKey[A]{
		key: internal.GenericKey[A]{},
		val: a,
	}
}

func Inject(assets ...syringe) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			for _, a := range assets {
				ctx = a.inject(ctx)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func GenericAssets[A any](a *A) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, internal.GenericKey[A]{}, a)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func GenericAssets2[A, B any](a *A, b *B) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, internal.GenericKey[A]{}, a)
			ctx = context.WithValue(ctx, internal.GenericKey[B]{}, b)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func GenericAssets3[A, B, C any](a *A, b *B, c *C) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, internal.GenericKey[A]{}, a)
			ctx = context.WithValue(ctx, internal.GenericKey[B]{}, b)
			ctx = context.WithValue(ctx, internal.GenericKey[C]{}, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func GenericAssets4[A, B, C, D any](a *A, b *B, c *C, d *D) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, internal.GenericKey[A]{}, a)
			ctx = context.WithValue(ctx, internal.GenericKey[B]{}, b)
			ctx = context.WithValue(ctx, internal.GenericKey[C]{}, c)
			ctx = context.WithValue(ctx, internal.GenericKey[D]{}, d)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
