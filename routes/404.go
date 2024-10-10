package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	. "github.com/Gardego5/htmdsl"
)

func Get404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	middleware.RenderPage(r, nil,
		components.Header{},
		components.Margins(
			"I don't exist yet. ",
			Code{"¯\\\\_(ツ)_/¯"},
		),
	)
}
