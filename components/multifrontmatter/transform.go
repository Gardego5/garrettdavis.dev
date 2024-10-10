package multifrontmatter

import (
	"encoding/json"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type (
	Page[T any] struct {
		Fm T
		Md string
	}

	transformer[T any] struct {
		Unmarshaller
		Splitter
	}
	Transformer[T any] interface {
		Transform(raw string) (Split[Page[T]], error)
	}
)

var _ Transformer[any] = (*transformer[any])(nil)

func NewTransformer[T any](format Unmarshaller, tag string) Transformer[T] {
	return transformer[T]{
		Unmarshaller: format,
		Splitter:     NewSplitter(format.LangTag(), tag),
	}
}

func (t transformer[T]) Transform(raw string) (Split[Page[T]], error) {
	unmarshalled := t.Split(raw)
	result := Split[Page[T]]{unmarshalled.Pre, make([]Page[T], 0, len(unmarshalled.Pages))}

	for _, page := range unmarshalled.Pages {
		var fm T
		if err := t.Unmarshal([]byte(page.Fm), &fm); err != nil {
			return result, err
		}
		result.Pages = append(result.Pages, Page[T]{Fm: fm, Md: page.Md})
	}

	return result, nil
}

type (
	Unmarshaller interface {
		Unmarshal([]byte, any) error
		LangTag() string
	}
	unmarshaller int
)

var _ Unmarshaller = (*unmarshaller)(nil)

const (
	Yaml unmarshaller = iota
	Toml
	Json
)

func (u unmarshaller) Unmarshal(b []byte, v any) error {
	switch u {
	case Yaml:
		return yaml.Unmarshal(b, v)
	case Toml:
		return toml.Unmarshal(b, v)
	case Json:
		return json.Unmarshal(b, v)
	default:
		panic("unreachable")
	}
}

func (u unmarshaller) LangTag() string {
	switch u {
	case Yaml:
		return "yaml"
	case Toml:
		return "toml"
	case Json:
		return "json"
	default:
		panic("unreachable")
	}
}
