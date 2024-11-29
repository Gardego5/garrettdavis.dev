package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	. "github.com/Gardego5/htmdsl"
)

func Get404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render.Page(w, r, nil,
		components.Header{},
		components.Margins{
			"I don't exist yet. ",
			Code{"¯\\_(ツ)_/¯"},
		},
	)
}
