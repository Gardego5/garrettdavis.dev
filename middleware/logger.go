package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/google/uuid"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		logger := global.Logger.With("requestId", id)

		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIdKey, id)
		ctx = context.WithValue(ctx, loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLogger(r *http.Request) *slog.Logger {
	if ptr := r.Context().Value(loggerKey); ptr != nil {
		return ptr.(*slog.Logger)
	} else {
		return nil
	}
}

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r)
		logger.Info("request received", "method", r.Method, "uri", r.RequestURI)
		lw := &loggingWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r)
		logger.Debug("request headers sent", "status", lw.statusCode, "headers", lw.Header())
	})
}
