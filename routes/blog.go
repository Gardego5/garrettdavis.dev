package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/blog"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

type Blog struct {
	blog *blog.Service
}

func NewBlog(
	blog *blog.Service,
) *Blog {
	return &Blog{
		blog: blog,
	}
}

func (h *Blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	post, found := h.blog.Posts()[r.PathValue("slug")]
	if !found {
		Get404(w, r)
		return
	}

	render.Page(w, r,
		Fragment{
			Title{post.Title},
			If(post.Description != "", Meta{"name": "description", "content": post.Description}),
			Style{PreEscaped(post.Css)},
		},
		components.Header{Title: post.Title},
		components.Margins(Div{Class("markdown"), PreEscaped(post.Content)}),
	)
}
