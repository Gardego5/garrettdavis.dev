package routes

import (
	_ "embed"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/resume"
	. "github.com/Gardego5/htmdsl"
)

type Resume struct{}

func NewResume() *Resume { return new(Resume) }

func (*Resume) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resume := access.Get[resume.Service](r.Context())

	render.Page(w, r,
		Fragment{
			Title{"Resume - Garrett Davis"},
			Meta{{"name", "description"}, {"content", "Garrett Davis' resume"}},
		},

		Div{Class("fixed top-0 right-0 flex gap-2 p-2 print:hidden"),
			Button{Attrs{{"x-data"}, {"x-on:click", "window.print()"}, {"title", "Print this page."}},
				Element("iconify-icon", Attrs{{"icon", "ph:printer"}, {"width", "36"}, {"height", "36"}}),
			},
		},

		Div{Class("mx-12"),
			Div{Class("grid grid-cols-[1fr_6fr] [&>*:nth-child(odd)]:text-right gap-4 print:text-black max-w-5xl m-auto justify-top pb-8"),
				Div{Class("col-start-2 flex py-4 items-center"),
					H1{
						Span{Class("text-4xl font-bold"), "Garrett Davis"},
						Br{},
						Span{Class("text-lg")},
					},
				},

				resume.DefaultResume(),
			},
		},
	)
}
