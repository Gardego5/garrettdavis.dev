package multifrontmatter_test

import (
	"reflect"
	"testing"

	. "github.com/Gardego5/garrettdavis.dev/utils/multifrontmatter"
)

type testTransform[T any] struct {
	transformer Transformer[T]
	in          string
	out         Split[Page[T]]
}

func (tt testTransform[T]) Test() func(*testing.T) {
	return func(t *testing.T) {
		out, _ := tt.transformer.Transform(tt.in)
		if out.Pre != tt.out.Pre {
			t.Fatalf("expected pre: %q, got %q", tt.out.Pre, out.Pre)
		}
		if len(out.Pages) != len(tt.out.Pages) {
			t.Fatalf("expected %d outputs, got %d", len(tt.out.Pages), len(out.Pages))
		}
		for i, o := range out.Pages {
			if !reflect.DeepEqual(o.Fm, tt.out.Pages[i].Fm) {
				t.Fatalf("expected frontmatter: %v, got %v", tt.out.Pages[i].Fm, o.Fm)
			}
			if o.Md != tt.out.Pages[i].Md {
				t.Fatalf("expected markdown: %q, got %q", tt.out.Pages[i].Md, o.Md)
			}
		}
	}
}

func TestTransform(t *testing.T) {
	type Titled struct {
		Title string `yaml:"title" toml:"title" json:"title"`
	}

	t.Run("yaml with frontmatter", testTransform[Titled]{
		transformer: NewTransformer[Titled](Yaml, "frontmatter"),
		in: multiline(
			"```yaml :frontmatter:",
			"title: hello world",
			"```",
			"",
			"this is a test",
		),
		out: Split[Page[Titled]]{
			Pages: []Page[Titled]{{
				Fm: Titled{Title: "hello world"},
				Md: "this is a test",
			}},
		},
	}.Test())

	t.Run("toml with frontmatter and pre", testTransform[Titled]{
		transformer: NewTransformer[Titled](Toml, "frontmatter"),
		in: multiline(
			"pre",
			"```toml :frontmatter:",
			`title = "hello world"`,
			"```",
			"",
			"this is a test",
		),
		out: Split[Page[Titled]]{
			Pre: "pre",
			Pages: []Page[Titled]{{
				Fm: Titled{Title: "hello world"},
				Md: "this is a test",
			}},
		},
	}.Test())

	t.Run("json with frontmatter and pre, with some invalid codeblocks", testTransform[Titled]{
		transformer: NewTransformer[Titled](Json, "frontmatter"),
		in: multiline(
			"```json :frontmatter:",
			"{\"title\": \"hello world\"}",
			"```",
			"",
			"this is a test",
			"```yaml",
			"title: hello world",
			"```",
		),
		out: Split[Page[Titled]]{
			Pages: []Page[Titled]{{
				Fm: Titled{Title: "hello world"},
				Md: multiline(
					"this is a test",
					"```yaml",
					"title: hello world",
					"```",
				),
			}},
		},
	}.Test())
}
