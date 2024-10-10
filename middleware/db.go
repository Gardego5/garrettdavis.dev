package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/Gardego5/goutils/env"
	"github.com/jmoiron/sqlx"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewDB() *sqlx.DB {
	env := env.MustLoad[struct {
		PrimaryUrl string `env:"TURSO_DATABASE_URL"`
		AuthToken  string `env:"TURSO_AUTH_TOKEN"`
	}]()

	driverName := "libsql"
	dataSourceName := env.PrimaryUrl + "?authToken=" + env.AuthToken

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		global.Logger.Info("Error opening database", "error", err)
		os.Exit(1)
	}

	return sqlx.NewDb(db, driverName)
}

func DB(db *sqlx.DB) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, dbKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDB(r *http.Request) *sqlx.DB {
	return r.Context().Value(dbKey).(*sqlx.DB)
}
