package access

import (
	"context"
	"log/slog"

	"github.com/Gardego5/garrettdavis.dev/resource/internal"
)

func Logger(c context.Context, scope string) *slog.Logger {
	return c.Value(internal.Logger).(*slog.Logger).With("scope", scope)
}
