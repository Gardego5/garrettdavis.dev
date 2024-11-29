package initialize

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func NewDB(databaseUrl, authToken string) *sqlx.DB {
	driverName := "libsql"
	dataSourceName := fmt.Sprintf("%s?authToken=%s", databaseUrl, authToken)

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		slog.Error("Error opening database", "error", err)
		os.Exit(1)
	}

	return sqlx.NewDb(db, driverName)
}

func NewRedis(url string) *redis.Client {
	opts, err := redis.ParseURL(url)
	if err != nil {
		slog.Error("Error parsing redis URL", "error", err)
		os.Exit(1)
	}

	return redis.NewClient(opts)
}
