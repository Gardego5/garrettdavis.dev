package middleware

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type md struct{ markdown goldmark.Markdown }

func (m md) Parse(src []byte) (io.Reader, parser.Context) {
	buf, ctx := &bytes.Buffer{}, parser.NewContext()
	m.markdown.Convert(src, buf, parser.WithContext(ctx))
	return buf, ctx
}

var MD = md{goldmark.New(
	goldmark.WithExtensions(
		&frontmatter.Extender{Mode: frontmatter.SetMetadata},

		highlighting.NewHighlighting(
			highlighting.WithGuessLanguage(true),
			highlighting.WithFormatOptions(
				chromahtml.WithClasses(true),
				chromahtml.WithAllClasses(true),
				chromahtml.WithLineNumbers(true),
			),
		),
	),
)}

func Markdown(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := context.WithValue(r.Context(), markdownKey, MD)
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func GetMarkdown(r *http.Request) md {
	return r.Context().Value(markdownKey).(md)
}

func GetFrontmatter(ctx parser.Context, dest any) error {
	data := frontmatter.Get(ctx)
	if data == nil {
		return errors.New("missing frontmatter")
	}
	return data.Decode(dest)
}
