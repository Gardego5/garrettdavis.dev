package routes

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/elliotchance/pie/v2"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

//go:embed blog/*
var blogfs embed.FS

type (
	BlogFrontmatter struct {
		Title       string    `yaml:"title"`
		Author      string    `yaml:"author"`
		Live        bool      `yaml:"live"`
		CreatedAt   time.Time `yaml:"createdAt"`
		UpdateAt    time.Time `yaml:"updatedAt"`
		Description string    `yaml:"description"`
	}
	BlogPost struct {
		Css, Content, Name string
		BlogFrontmatter
	}
	BlogPostSummary struct{ *BlogPost }
)

var postsMap = sync.OnceValue(func() (postsMap map[string]BlogPost) {
	postsMap = make(map[string]BlogPost)

	dirents, err := blogfs.ReadDir("blog")
	if err != nil {
		global.Logger.Error("error reading blog directory", "error", err)
		os.Exit(1)
	}

	for _, dirent := range dirents {
		name := dirent.Name()
		if dirent.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		file, err := blogfs.ReadFile(fmt.Sprintf("blog/%s", name))
		if err != nil {
			global.Logger.Error("error reading blog file", "error", err)
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
		fm := BlogFrontmatter{}
		content, ctx := &bytes.Buffer{}, parser.NewContext()
		err = md.Convert(file, content, parser.WithContext(ctx))
		if err != nil {
			global.Logger.Error("error converting markdown", "error", err)
			os.Exit(1)
		}

		err = middleware.GetFrontmatter(ctx, &fm)
		if err != nil {
			global.Logger.Error("error parsing frontmatter", "error", err)
			os.Exit(1)
		}

		postsMap[name] = BlogPost{
			Css: css.String(), Content: content.String(), Name: name,
			BlogFrontmatter: fm,
		}
	}

	return
})

var postsList = sync.OnceValue(func() (postsList []BlogPost) {
	postsList = pie.Values(postsMap())

	slices.SortStableFunc(postsList, func(a, z BlogPost) int {
		return int(z.CreatedAt.Unix() - a.CreatedAt.Unix())
	})

	return
})

var livePostsList = sync.OnceValue(func() (livePostsList []BlogPost) {
	livePostsList = pie.Filter(postsList(), func(post BlogPost) bool { return post.Live })

	return
})

func (p BlogPost) Summary() BlogPostSummary { return BlogPostSummary{&p} }

func (p BlogPostSummary) Render() RenderedHTML {
	return Div{Class("flex flex-col gap-2 my-4"),
		H2{Class("text-xl"),
			A{Attrs{{"href", fmt.Sprintf("/blog/%s", p.Name)}}, p.Title},
		},
		P{Class("text-gray-500 text-base"), p.Description},
	}.Render()
}

func GetBlogPage(w http.ResponseWriter, r *http.Request) {
	post, found := postsMap()[r.PathValue("slug")]
	if !found {
		Get404(w, r)
		return
	}

	middleware.RenderPage(r,
		Fragment{
			Title{post.Title},
			If(post.Description != "", Meta{{"name", "description"}, {"content", post.Description}}),
			Style{PreEscaped(post.Css)},
		},
		components.Header{Title: post.Title},
		components.Margins(Div{Class("markdown"), PreEscaped(post.Content)}),
	)
}
