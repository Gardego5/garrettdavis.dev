package routes

import (
	_ "embed"
	"fmt"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/presentations"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

type Presentations struct {
	presentations *presentations.Service
}

var (
	//go:embed presentation.go.css
	presentationCSS PreEscaped
	//go:embed presentation.go.js
	presentationJS PreEscaped
)

func NewPresentations(
	presentations *presentations.Service,
) *Presentations {
	return &Presentations{
		presentations: presentations,
	}
}

func (p *Presentations) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pres, found := p.presentations.Presentations()[r.PathValue("slug")]
	if !found {
		Get404(w, r)
		return
	}

	slidesLen := len(pres.Slides)
	slides := make(Fragment, 0, slidesLen)
	for idx, slide := range pres.Slides {
		slides = append(slides,
			Section{Id(fmt.Sprintf("%d", idx+1)), Class("relative", "h-screen"),
				Div{Class("p-20 inset-4 absolute border border-slate-500 slide grid grid-rows-[auto_minmax(0,1fr)] gap-16"),
					Div{Class("grid gap-8 presentation", slide.SlideFrontmatter.Class),
						slide.Content,
					},
				},

				Button{Class("absolute top-8 left-8 text-slate-500 text-xs italic p-1"),
					Attrs{"@click": "console.log($refs.modal), $refs.modal.closeModal()"},
					idx + 1, "/", slidesLen,
				},

				If(idx > 0, Button{Class("w-6 h-6 text-slate-500 absolute right-8 top-8"),
					Attrs{"@click": "decrement"},
					Element("iconify-icon", Attrs{"icon": "ph:caret-up", "width": 24, "height": 24}),
				}),

				If(idx+1 < slidesLen, Button{
					Class("w-6 h-6 text-slate-500 absolute right-8 bottom-8"),
					Attrs{"@click": "increment"},
					Element("iconify-icon", Attrs{"icon": "ph:caret-down", "width": 24, "height": 24}),
				}),
			},
		)
	}

	render.Page(w, r,
		Fragment{
			Style{presentationCSS},
			Script{presentationJS},
		},
		Main{Attrs{
			"x-data":                             fmt.Sprintf(`{slide:0,modal:false,max:%d,increment(){this.slide<this.max&&this.slide++},decrement(){this.slide>1&&this.slide--}}`, slidesLen),
			"@keydown.arrow-down.window.prevent": "increment",
			"@keydown.arrow-up.window.prevent":   "decrement",
			"@keydown.space.window.prevent":      "increment",
			"@keydown.home.window.prevent":       "slide=1",
			"@keydown.end.window.prevent":        fmt.Sprintf("slide=%d", slidesLen),
			"x-effect":                           "window.document.getElementById(slide.toString())?.scrollIntoView()",
		},
			Div{Class("transition-all delay-500 duration-1000"),
				Attrs{
					"style":  "min-height: 60vh",
					"x-init": "$el.style.minHeight = '100vh'",
				},
				components.Header{Variant: components.HeaderStatic},
				components.Margins(Div{Class("markdown"), pres.Content}),
			},

			slides,

			Dialog{Class("bg-black/50 p-3 rounded-lg"),
				Attrs{"x-ref": "modal"},
				Form{Attrs{"method": "dialog"},
					Input{
						"class":        "bg-transparent text-slate-300 text-center text-lg",
						"@beforeinput": fmt.Sprintf("(NUMBERS_ONLY.test($event.data) && !($event.target.value+$event.data) || inRange(+$event.target.value+$event.data, 0, %d)) || $event.preventDefault()", slidesLen),
						"autofocus":    nil,
						"type":         "text",
						"x-model":      "slide",
					},
				},
			},
		},
	)
}
