package middleware

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/htmdsl"
)

func RenderPage(r *http.Request, head html.HTML, body ...any) {
	w := r.Context().Value(writerrefKey).(http.ResponseWriter)
	html.Render(w, components.Page{
		Head: head, Body: html.Fragment{body},
		HxUpdate: r.Header.Get("hx-boosted") == "true",
	})
}
