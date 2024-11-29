package routes

import (
	"context"
	_ "embed"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
	"github.com/elliotchance/pie/v2"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type (
	ResumeItem struct {
		Annotation string   `json:"annotation" yaml:"annotation"`
		Title      string   `json:"title" yaml:"title"`
		Subtitle   string   `json:"subtitle" yaml:"subtitle"`
		Block      string   `json:"block,omitempty" yaml:"block,omitempty" validate:"required_without=List,excluded_if=List"`
		List       []string `json:"list,omitempty" yaml:"list,omitempty" validate:"required_without=Block,excluded_if=Block"`
	}
	ResumeSection struct {
		Title string       `json:"title" yaml:"title"`
		Items []ResumeItem `json:"items" yaml:"items"`
	}
	ResumeDoc struct {
		Preface  string          `json:"preface" yaml:"preface"`
		Sections []ResumeSection `json:"sections" yaml:"sections"`
	}
)

func (i ResumeItem) Render(ctx context.Context) RenderedHTML {
	return Div{Class("pb-3"),
		H3{Class("text-blue-300 text-xl"),
			i.Title,

			If(i.Annotation != "", Span{Class("inline-block float-right"),
				i.Annotation,
			}),
		},

		If(i.Subtitle != "", H4{Class("text-slate-300 text-sm pl-2 -mt-1 mb-1 italic print:text-xs"),
			i.Subtitle,
		}),

		If(i.Block == "",
			Ul{Class("flex gap-x-3 gap-y-2 flex-wrap text-sm"), pie.Map(i.List, func(s string) HTML { return Li{s} })}).
			Else(P{Class("text-sm"), i.Block}),
	}.Render(ctx)
}

func (s ResumeSection) Render(context.Context) RenderedHTML {
	return Fragment{
		H2{Class("font-mono leading-7 tracking-lighter"), s.Title},
		Div{pie.Map(s.Items, func(i ResumeItem) any { return i })},
	}
}

func (d ResumeDoc) Render(context.Context) RenderedHTML {
	return Fragment{
		P{Class("col-start-2"), d.Preface},

		H2{Class("font-mono leading-7 tracking-lighter col-start-1"), "Contact"},
		Ul{Class("col-start-2"),
			Li{"Email: ", A{Attrs{
				{"x-data", `{user:'contact',domain:window.location.hostname}`},
				{"x-init", `$el.href = 'mailto:' + $data.user + '@' + $data.domain`},
			},
				"contact [at] ", Span{Attrs{{"x-text", "domain"}}},
			}},

			Li{"Location: ", "Hillsboro, Oregon"},
		},

		d.Sections,
	}
}

//go:embed resume.yaml
var yamlDoc []byte

type Resume struct{}

func NewResume() *Resume { return new(Resume) }

func (*Resume) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := access.Logger(r.Context(), "Resume")

	var doc ResumeDoc
	if err := yaml.Unmarshal(yamlDoc, &doc); err != nil {
		logger.Error("Failed to unmarshal yaml", "err", err)
		return
	}

	logger.Info("Validate resume",
		"err", access.Get[validator.Validate](r.Context()).Struct(doc))

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
			Div{Class("grid grid-cols-[1fr_6fr] [&>*:nth-child(odd)]:text-right gap-4 print:text-black max-w-5xl m-auto justify-top"),
				Div{Class("col-start-2 flex py-4 items-center"),
					H1{
						Span{Class("text-4xl font-mono font-bold"), "Garrett Davis"},
						Br{},
						Span{Class("text-lg")},
					},
				},

				doc,
			},
		},
	)
}
