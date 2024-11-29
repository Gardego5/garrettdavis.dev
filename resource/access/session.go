package access

import (
	"context"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/resource/internal"
	"github.com/Gardego5/garrettdavis.dev/utils/cookie"
)

func Session(c context.Context) string {
	logger := Logger(c, "Session")
	w := c.Value(internal.WriterRef).(http.ResponseWriter)

	for _, line := range w.Header()["Set-Cookie"] {
		c, err := http.ParseSetCookie(line)
		if err != nil {
			continue
		} else if c.Name == cookie.Session.Name {
			return c.Value
		}
	}

	session, ok := c.Value(internal.Session).(string)
	if !ok {
		logger.Warn("session not found in context")
	}

	return session
}
