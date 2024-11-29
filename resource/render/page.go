package render

import (
	"context"
	"net/http"

	"github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/util"
	"github.com/elliotchance/pie/v2"
)

func Page(w http.ResponseWriter, r *http.Request, head html.HTML, body ...any) {
	html.RenderContext(w, r.Context(), page{
		Head: head, Body: html.Fragment{body},
		HxUpdate: r.Header.Get("hx-boosted") == "true",
	})
}

type page struct {
	Head, Body html.HTML
	HxUpdate   bool
}

var _ html.HTML = (*page)(nil)

func (page page) Render(context.Context) html.RenderedHTML {
	return html.Fragment{html.DOCTYPE, html.Html{html.Attrs{{"lang", "en"}},
		html.Head{
			page.Head,

			util.If(!page.HxUpdate,
				// Meta Tags
				html.Meta{{"charset", "utf-8"}, {"hx-preserve", "true"}},
				html.Meta{{"name", "viewport"}, {"content", "width=device-width, initial-scale=1"}, {"hx-preserve", "true"}},

				// Tailwind CSS
				html.Link{{"rel", "stylesheet"}, {"href", "/static/css/style.css"}, {"hx-preserve", "true"}},

				// Local Third-Party Scripts
				html.Script{html.Attrs{{"src", "/static/3p/js/htmx.2.0.1.min.js"}, {"hx-preserve", "true"}}},
				html.Script{html.Attrs{{"src", "/static/3p/js/htmx-ext-multi-swap.2.0.0.js"}, {"hx-preserve", "true"}}},
				html.Script{html.Attrs{{"src", "/static/3p/js/htmx-ext-preload.2.0.1.js"}, {"hx-preserve", "true"}}},
				html.Script{html.Attrs{{"src", "/static/3p/js/htmx-ext-response-targets.2.0.0.js"}, {"hx-preserve", "true"}}},
				html.Script{html.Attrs{{"src", "/static/3p/js/htmx-ext-head-support.2.0.1.js"}, {"hx-preserve", "true"}}},
				pie.Map([]string{
					"/static/3p/js/iconify-icon.2.1.0.min.js",
					"/static/3p/js/alpinejs-anchor.3.14.1.min.js",
					"/static/3p/js/alpinejs-morph.3.14.1.min.js",
					"/static/3p/js/alpinejs.3.14.1.min.js",
				}, func(src string) html.RenderedHTML {
					// Even though we want to defer the execution of these
					// scripts, we don't want to delay it's loading.
					return html.Fragment{
						html.Link{{"rel", "preload"}, {"as", "script"}, {"href", src}, {"hx-preserve", "true"}},
						html.Script{html.Attrs{{"defer"}, {"src", src}, {"hx-preserve", "true"}}},
					}
				}),

				html.Link{{"rel", "stylesheet"}, {"href", "/static/fonts/truetype/IosevkaGarrettDavisDev.css"}, {"hx-preserve", "true"}},
			),
		},
		html.Body{html.Attrs{{"class", "box-border bg-zinc-950 text-zinc-50"}, {"hx-boost", "true"}, {"hx-ext", "response-targets,head-support"}},
			html.Div{html.Attrs{{"class", "bg-zinc-900 min-h-[100vh]"}},
				page.Body,
			},
		},
	}}
}
