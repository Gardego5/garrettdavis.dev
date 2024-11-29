package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/resource/internal"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/utils/cookie"
	"github.com/Gardego5/garrettdavis.dev/utils/mux"
	html "github.com/Gardego5/htmdsl"
	"github.com/google/uuid"
)

func logger(ctx context.Context) *slog.Logger {
	return ctx.Value(internal.Logger).(*slog.Logger)
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LoggerAndSessions(logger *slog.Logger, logRequests bool) mux.Middleware {
	return mux.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			method, path, requestId := r.Method, r.RequestURI, uuid.NewString()
			logger := logger.With("requestId", requestId)

			// initialize session cookie if it's missing
			var session string
			if Cookie, err := r.Cookie(cookie.Session.Name); err == http.ErrNoCookie {
				logger.Debug("setting new session cookie", "value", session)
				session = uuid.NewString()
				NewCookie := cookie.Session
				NewCookie.Value = session
				http.SetCookie(w, &NewCookie)
			} else if err != nil {
				logger.Error("error getting session cookie", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				render.Page(w, r, nil, html.P{"Something went wrong."})
				return
			} else {
				logger.Debug("got existing session cookie", "value", Cookie.Value)
				session = Cookie.Value
			}

			logger = logger.With("session", session)

			ctx = context.WithValue(ctx, internal.Logger, logger)
			ctx = context.WithValue(ctx, internal.RequestId, requestId)
			ctx = context.WithValue(ctx, internal.RequestRef, r)
			ctx = context.WithValue(ctx, internal.RouterMethod, method)
			ctx = context.WithValue(ctx, internal.RouterPath, path)
			ctx = context.WithValue(ctx, internal.Session, session)
			ctx = context.WithValue(ctx, internal.WriterRef, w)

			if logRequests {
				logger.Info("request received", "method", r.Method, "uri", r.RequestURI)
			}

			lw := &loggingWriter{ResponseWriter: w}
			next.ServeHTTP(lw, r.WithContext(ctx))

			if logRequests {
				logger.Debug("request headers sent", "status", lw.statusCode, "headers", lw.Header())
			}
		})
	})
}
