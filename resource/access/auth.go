//go:build pendingsettingupaccess

package access

import (
	"context"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/schema"
	"github.com/Gardego5/garrettdavis.dev/utils/data"
	"github.com/Gardego5/garrettdavis.dev/utils/middleware/internal"
	"github.com/casbin/casbin/v2"
	"github.com/google/go-github/v66/github"
)

func Enforcer(ctx context.Context) *casbin.Enforcer {
	return ctx.Value(internal.Enforcer).(*casbin.Enforcer)
}

func Subject(ctx context.Context) (*schema.Subject, error) {
	return Cache[schema.Subject](ctx).Get(ctx, Session(ctx))
}

func CurrentUser(ctx context.Context) (*github.User, error) {
	session := Session(ctx)

	user, err := data.GetCurrentUser(
		ctx, session,
		Cache[schema.GHAccessToken](ctx),
		Cache[github.User](ctx))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Enforce(ctx context.Context, path, method string) (bool, error) {
	if sub, err := Subject(ctx); err != nil {
		return false, err
	} else {
		return Enforcer(ctx).Enforce(sub, path, method)
	}
}

func RequireAuth() mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := Logger(ctx, "RequireAuth")

			if has, err := Enforce(ctx, Path(ctx), r.Method); err != nil {
				logger.Warn("Error enforcing policy", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
			} else if !has {
				logger.Debug("User does not have permission")
				w.WriteHeader(http.StatusForbidden)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	})
}
