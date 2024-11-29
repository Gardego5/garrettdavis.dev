package presentations

import (
	"bytes"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/Gardego5/garrettdavis.dev/utils/multifrontmatter"
	. "github.com/Gardego5/htmdsl"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed data/*
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

type Service struct {
	presentations map[string]Presentation

	once sync.Once
}

func New() *Service { return new(Service) }

func (s *Service) Presentations() map[string]Presentation {
	s.once.Do(s.initialize)
	return s.presentations
}

func (s *Service) initialize() {
	s.presentations = make(map[string]Presentation)

	dirents, err := presentationfs.ReadDir("presentation")
	if err != nil {
		slog.Error("error reading blog directory", "error", err)
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
			slog.Error("error reading blog file", "error", err)
			os.Exit(1)
		}

		// Transform the presentation file into frontmatter and content
		presentationS, err := presentationT.Transform(string(file))
		if err != nil {
			slog.Error("error transforming presentation", "error", err)
			os.Exit(1)
		} else if len(presentationS.Pages) != 1 {
			slog.Error("error transforming presentation",
				"error", fmt.Errorf("expected 1 page, got %d", len(presentationS.Pages)))
			os.Exit(1)
		} else if len(presentationS.Pre) != 0 {
			slog.Error("error transforming presentation",
				"error", fmt.Errorf("expected len(pre) == 0, got %d", len(presentationS.Pre)))
			os.Exit(1)
		}

		presPage := presentationS.Pages[0]
		presentation := Presentation{
			PresentationFrontmatter: presPage.Fm,
			Name:                    strings.TrimSuffix(name, ".md"),
		}

		slides, err := slideT.Transform(presPage.Md)
		if err != nil {
			slog.Error("error transforming slides", "error", err)
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
		err = md.Convert([]byte(slides.Pre), content)
		if err != nil {
			slog.Error("error converting markdown", "error", err)
			os.Exit(1)
		}
		presentation.Content = PreEscaped(content.String())
		presentation.Slides = make([]Slide, 0, len(slides.Pages))

		for _, sld := range slides.Pages {
			content := &bytes.Buffer{}
			err = md.Convert([]byte(sld.Md), content)
			if err != nil {
				slog.Error("error converting markdown", "error", err)
				os.Exit(1)
			}
			presentation.Slides = append(presentation.Slides, Slide{
				SlideFrontmatter: sld.Fm, Content: PreEscaped(content.String()),
			})
		}

		presentation.Css = PreEscaped(css.String())
		s.presentations[presentation.Name] = presentation
	}
}
