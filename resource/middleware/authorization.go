package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	"github.com/Gardego5/garrettdavis.dev/utils/mux"
	"github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
)

func Authorization(
	logger *slog.Logger,
	enforcer *casbin.Enforcer,
	currentuser *currentuser.Service,
) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		logger := logger.With("scope", "middleware.Authorization")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			session := access.Session(ctx)
			user, err := currentuser.GetUserBySession(ctx, session)
			if err == redis.Nil {
				logger.WarnContext(ctx, "session not found... the user probably hasn't signed in.")
				w.WriteHeader(http.StatusUnauthorized)
				return
			} else if err != nil {
				logger.ErrorContext(ctx, "error getting user by session", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			sub := model.Subject{User: user.GetLogin()}

			if ok, err := enforcer.Enforce(sub, r.URL.Path, r.Method); err != nil {
				logger.ErrorContext(ctx, "error enforcing policy", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else if !ok {
				logger.WarnContext(ctx, "unauthorized access", "subject", sub)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
			return
		})
	})
}
