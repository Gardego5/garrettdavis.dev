package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func Must[T any](val T, err error) T {
	if err != nil {
		slog.Error(
			fmt.Sprintf("called Must[%T](%v, %s) but err != nil", val, val, err),
			"error", err)
		os.Exit(1)
	}
	return val
}
