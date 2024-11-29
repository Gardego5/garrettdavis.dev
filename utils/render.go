package utils

import (
	"errors"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

func Frontmatter(ctx parser.Context, dest any) error {
	data := frontmatter.Get(ctx)
	if data == nil {
		return errors.New("missing frontmatter")
	}
	return data.Decode(dest)
}
