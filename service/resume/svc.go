package resume

import (
	"context"
	_ "embed"
	"log/slog"
	"os"
	"sync"

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
		H2{Class("leading-7 tracking-lighter"), s.Title},
		Div{pie.Map(s.Items, func(i ResumeItem) any { return i })},
	}
}

func (d ResumeDoc) Render(context.Context) RenderedHTML {
	return Fragment{
		P{Class("col-start-2"), d.Preface},

		H2{Class("leading-7 tracking-lighter col-start-1"), "Contact"},
		Ul{Class("col-start-2"),
			Li{"Email: ", A{Attrs{
				"x-data": `{user:'contact',domain:window.location.hostname}`,
				"x-init": `$el.href = 'mailto:' + $data.user + '@' + $data.domain`,
			},
				"contact [at] ", Span{Attrs{"x-text": "domain"}},
			}},

			Li{"Location: ", "Hillsboro, Oregon"},
		},

		d.Sections,
	}
}

//go:embed data/default-resume.yaml
var defaultResume []byte

type Service struct {
	defaultResume ResumeDoc

	validate *validator.Validate

	once sync.Once
}

func New(validate *validator.Validate) *Service { return &Service{validate: validate} }

func (svc *Service) DefaultResume() *ResumeDoc { svc.once.Do(svc.init); return &svc.defaultResume }

func (svc *Service) init() {
	if err := yaml.Unmarshal(defaultResume, &svc.defaultResume); err != nil {
		slog.Error("failed to unmarshal yaml", "err", err)
		os.Exit(1)
	}

	if err := svc.validate.Struct(svc.defaultResume); err != nil {
		slog.Error("failed to validate default resume", "err", err)
		os.Exit(1)
	}
}
