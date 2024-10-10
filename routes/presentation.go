package routes

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/components/multifrontmatter"
	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/k0kubun/pp"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed presentation/*
var presentationfs embed.FS

type (
	Presentation struct {
		PresentationFrontmatter
		Css, Content PreEscaped
		Name         string
		Slides       []Slide
	}
	PresentationFrontmatter struct {
		Live bool `yaml:"live"`
	}
	Slide struct {
		SlideFrontmatter
		Content PreEscaped
	}
	SlideFrontmatter struct {
		Class string `yaml:"class"`
	}
)

var presentationMap = sync.OnceValue(func() (presentationMap map[string]Presentation) {
	presentationMap = make(map[string]Presentation)

	dirents, err := presentationfs.ReadDir("presentation")
	if err != nil {
		global.Logger.Error("error reading blog directory", "error", err)
		os.Exit(1)
	}

	presentationT := multifrontmatter.NewTransformer[PresentationFrontmatter](
		multifrontmatter.Yaml, "presentation")
	slideT := multifrontmatter.NewTransformer[SlideFrontmatter](
		multifrontmatter.Yaml, "slide")

	for _, dirent := range dirents {
		name := dirent.Name()
		if dirent.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		file, err := presentationfs.ReadFile(fmt.Sprintf("presentation/%s", name))
		if err != nil {
			global.Logger.Error("error reading blog file", "error", err)
			os.Exit(1)
		}

		// Transform the presentation file into frontmatter and content
		presentationS, err := presentationT.Transform(string(file))
		if err != nil {
			global.Logger.Error("error transforming presentation", "error", err)
			os.Exit(1)
		} else if len(presentationS.Pages) != 1 {
			global.Logger.Error("error transforming presentation",
				"error", fmt.Errorf("expected 1 page, got %d", len(presentationS.Pages)))
			os.Exit(1)
		} else if len(presentationS.Pre) != 0 {
			global.Logger.Error("error transforming presentation",
				"error", fmt.Errorf("expected len(pre) == 0, got %d", len(presentationS.Pre)))
			os.Exit(1)
		}

		presPage := presentationS.Pages[0]
		presentation := Presentation{
			PresentationFrontmatter: presPage.Fm,
			Name:                    strings.TrimSuffix(name, ".md"),
		}

		slideS, err := slideT.Transform(presPage.Md)
		if err != nil {
			global.Logger.Error("error transforming slides", "error", err)
			os.Exit(1)
		}

		css := &bytes.Buffer{}
		md := goldmark.New(
			goldmark.WithRendererOptions(
				html.WithXHTML(),
				html.WithUnsafe(),
			),
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithCSSWriter(css),
					highlighting.WithGuessLanguage(true),
					highlighting.WithStyle("native"),
					highlighting.WithFormatOptions(
						chromahtml.WithClasses(true),
						chromahtml.WithAllClasses(true),
						chromahtml.WithLineNumbers(true),
					),
				),
			),
		)

		content := &bytes.Buffer{}
		err = md.Convert([]byte(slideS.Pre), content)
		if err != nil {
			global.Logger.Error("error converting markdown", "error", err)
			os.Exit(1)
		}
		presentation.Content = PreEscaped(content.String())
		presentation.Slides = make([]Slide, 0, len(slideS.Pages))

		for _, slide := range slideS.Pages {
			content := &bytes.Buffer{}
			err = md.Convert([]byte(slide.Md), content)
			if err != nil {
				global.Logger.Error("error converting markdown", "error", err)
				os.Exit(1)
			}
			presentation.Slides = append(presentation.Slides, Slide{
				SlideFrontmatter: slide.Fm, Content: PreEscaped(content.String()),
			})
		}

		presentation.Css = PreEscaped(css.String())
		presentationMap[presentation.Name] = presentation
	}

	pp.Print(presentationMap)

	return
})

var (
	//go:embed presentation.go.css
	presentationCSS PreEscaped
	//go:embed presentation.go.js
	presentationJS PreEscaped
)

func GetPresentation(w http.ResponseWriter, r *http.Request) {
	pres, found := presentationMap()[r.PathValue("slug")]
	if !found {
		Get404(w, r)
		return
	}

	slidesLen := len(pres.Slides)
	slides := make(Fragment, 0, slidesLen)
	for idx, slide := range pres.Slides {
		slides = append(slides,
			Section{Attrs{Class("relative", "h-screen"), {"id", fmt.Sprintf("%d", idx+1)}},
				Div{Class("p-20 inset-4 absolute border border-slate-500 slide grid grid-rows-[auto_minmax(0,1fr)] gap-16"),
					Div{Class("grid gap-8 presentation", slide.SlideFrontmatter.Class),
						slide.Content,
					},
				},

				Button{Attrs{
					Class("absolute top-8 left-8 text-slate-500 text-xs italic p-1"),
					{"@click", "console.log($refs.modal), $refs.modal.closeModal()"},
				},
					idx + 1, "/", slidesLen,
				},

				If(idx > 0, Button{Attrs{
					Class("w-6 h-6 text-slate-500 absolute right-8 top-8"),
					{"@click", "decrement"},
				},
					Element("iconify-icon", Attrs{{"icon", "ph:caret-up"}, {"width", "24"}, {"height", "24"}}),
				}),

				If(idx+1 < slidesLen, Button{Attrs{
					Class("w-6 h-6 text-slate-500 absolute right-8 bottom-8"),
					{"@click", "increment"},
				},
					Element("iconify-icon", Attrs{{"icon", "ph:caret-down"}, {"width", "24"}, {"height", "24"}}),
				}),
			},
		)
	}

	middleware.RenderPage(r,
		Fragment{
			Style{presentationCSS},
			Script{presentationJS},
			//Title{pres.Title},
		},
		Main{Attrs{
			{"x-data", fmt.Sprintf(`{slide:0,modal:false,max:%d,increment(){this.slide<this.max&&this.slide++},decrement(){this.slide>1&&this.slide--}}`, slidesLen)},
			{"@keydown.arrow-down.window.prevent", "increment"},
			{"@keydown.arrow-up.window.prevent", "decrement"},
			{"@keydown.space.window.prevent", "increment"},
			{"@keydown.home.window.prevent", "slide=1"},
			{"@keydown.end.window.prevent", fmt.Sprintf("slide=%d", slidesLen)},
			{"x-effect", "window.document.getElementById(slide.toString())?.scrollIntoView()"},
		},
			Div{Attrs{Class("transition-all delay-500 duration-1000"), {"style", "min-height: 60vh"}, {"x-init", "$el.style.minHeight = '100vh'"}},
				components.Header{Variant: components.HeaderStatic},
				components.Margins(Div{Class("markdown"), pres.Content}),
			},

			slides,

			Dialog{Attrs{Class("bg-black/50 p-3 rounded-lg"), {"x-ref", "modal"}},
				Form{Attrs{{"method", "dialog"}},
					Input{
						Class("bg-transparent text-slate-300 text-center text-lg"),
						{"@beforeinput", fmt.Sprintf("(NUMBERS_ONLY.test($event.data) && !($event.target.value+$event.data) || inRange(+$event.target.value+$event.data, 0, %d)) || $event.preventDefault()", slidesLen)},
						{"autofocus"}, {"type", "text"}, {"x-model", "slide"},
					},
				},
			},
		},
	)
}
