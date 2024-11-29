package blog

import (
	"bytes"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Gardego5/garrettdavis.dev/utils"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/elliotchance/pie/v2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type (
	Frontmatter struct {
		Title       string    `yaml:"title"`
		Author      string    `yaml:"author"`
		Live        bool      `yaml:"live"`
		CreatedAt   time.Time `yaml:"createdAt"`
		UpdateAt    time.Time `yaml:"updatedAt"`
		Description string    `yaml:"description"`
	}
	Post struct {
		Css, Content, Name string
		Frontmatter
	}
)

//go:embed data/*
var data embed.FS

type Service struct {
	posts      map[string]Post
	list, live []Post

	once sync.Once
}

func New() *Service { return new(Service) }

func (svc *Service) Posts() map[string]Post { svc.once.Do(svc.init); return svc.posts }
func (svc *Service) List() []Post           { svc.once.Do(svc.init); return svc.list }
func (svc *Service) Live() []Post           { svc.once.Do(svc.init); return svc.live }

func (svc *Service) init() {
	svc.posts = make(map[string]Post)

	dirents, err := data.ReadDir("data")
	if err != nil {
		slog.Error("error reading blog directory", "error", err)
		os.Exit(1)
	}

	for _, dirent := range dirents {
		name := dirent.Name()
		if dirent.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		file, err := data.ReadFile(fmt.Sprintf("data/%s", name))
		if err != nil {
			slog.Error("error reading blog file", "error", err)
			os.Exit(1)
		}

		css := &bytes.Buffer{}
		md := goldmark.New(
			goldmark.WithExtensions(
				&frontmatter.Extender{Mode: frontmatter.SetMetadata},
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

		name = strings.TrimSuffix(name, ".md")
		fm := Frontmatter{}
		content, ctx := &bytes.Buffer{}, parser.NewContext()
		err = md.Convert(file, content, parser.WithContext(ctx))
		if err != nil {
			slog.Error("error converting markdown", "error", err)
			os.Exit(1)
		}

		err = utils.Frontmatter(ctx, &fm)
		if err != nil {
			slog.Error("error parsing frontmatter", "error", err)
			os.Exit(1)
		}

		svc.posts[name] = Post{
			Css: css.String(), Content: content.String(), Name: name,
			Frontmatter: fm,
		}
	}

	svc.list = pie.Values(svc.posts)
	slices.SortStableFunc(svc.list, func(a, z Post) int {
		return int(z.CreatedAt.Unix() - a.CreatedAt.Unix())
	})

	svc.live = pie.Filter(svc.list, func(post Post) bool { return post.Live })
}
