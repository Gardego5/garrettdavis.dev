package global

import (
	"log/slog"
	"os"

	"github.com/Gardego5/goutils/env"
)

var (
	LogLevel = &slog.LevelVar{}
	Logger   = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: LogLevel}))
)

func init() {
	env := env.MustLoad[struct {
		LogLevel string `env:"LOG_LEVEL=INFO" validate:"oneof=DEBUG INFO WARN ERROR"`
	}]()

	if err := Validate.Struct(env); err != nil {
		Logger.Warn("invalid log level environment variable config", "error", err)
	} else {
		LogLevel.UnmarshalText([]byte(env.LogLevel))
	}
}
