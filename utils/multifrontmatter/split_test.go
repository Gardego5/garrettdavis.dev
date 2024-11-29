package multifrontmatter_test

import (
	"strings"
	"testing"

	. "github.com/Gardego5/garrettdavis.dev/utils/multifrontmatter"
)

func multiline(s ...string) string { return strings.Join(s, "\n") }

type testSplit struct {
	splitter Splitter
	name, in string
	out      Split[Data]
}

func TestSplit(t *testing.T) {
	tests := []testSplit{
		{
			splitter: NewSplitter("yaml", "frontmatter"),
			name:     "empty",
			in:       "",
			out:      Split[Data]{},
		},
		{
			splitter: NewSplitter("yaml", "frontmatter"),
			name:     "no frontmatter",
			in: multiline(
				"# hello world",
				"",
				"this is a test",
			),
			out: Split[Data]{Pre: multiline(
				"# hello world",
				"",
				"this is a test",
			)},
		},
		{
			splitter: NewSplitter("yaml", "frontmatter"),
			name:     "single frontmatter",
			in: multiline(
				"```yaml :frontmatter:",
				"title: hello world",
				"```",
				"",
				"this is a test",
			),
			out: Split[Data]{Pages: []Data{
				{
					Fm: multiline(
						"title: hello world",
					),
					Md: multiline(
						"this is a test",
					),
				},
			}},
		},
		{
			splitter: NewSplitter("yaml", "frontmatter"),
			name:     "multiple frontmatter",
			in: multiline(
				"```yaml :frontmatter:",
				"title: hello world",
				"```",
				"",
				"this is a test",
				"",
				"```yaml :frontmatter:",
				"title: hello universe",
				"```",
				"",
				"this is another test",
			),
			out: Split[Data]{Pages: []Data{
				{
					Fm: multiline(
						"title: hello world",
					),
					Md: multiline(
						"this is a test",
					),
				},
				{
					Fm: multiline(
						"title: hello universe",
					),
					Md: multiline(
						"this is another test",
					),
				},
			}},
		},
		{
			splitter: NewSplitter("toml", "slide"),
			name:     "toml slide",
			in: multiline(
				"this is prefixed text",
				"",
				"```toml :slide:",
				"layout = 1",
				"```",
				"",
				"this is a test",
				"",
				"```toml :slide:",
				"layout = 2",
				"```",
				"",
				"this is another test",
			),
			out: Split[Data]{
				Pre: "this is prefixed text",
				Pages: []Data{
					{
						Fm: multiline(
							"layout = 1",
						),
						Md: multiline(
							"this is a test",
						),
					},
					{
						Fm: multiline(
							"layout = 2",
						),
						Md: multiline(
							"this is another test",
						),
					},
				},
			},
		},
		{
			splitter: NewSplitter("json", "faq"),
			name:     "json faq long blocks",
			in: multiline(
				"## Frequently Asked Questions",
				"",
				"```json :faq:",
				"{",
				`  "question": "What is the meaning of life?"`,
				"}",
				"```",
				"",
				"The meaning of life is 42.",
				"",
				"Something in the way...",
				"",
			),
			out: Split[Data]{
				Pre: multiline(
					"## Frequently Asked Questions",
				),
				Pages: []Data{
					{
						Fm: multiline(
							"{",
							`  "question": "What is the meaning of life?"`,
							"}",
						),
						Md: multiline(
							"The meaning of life is 42.",
							"",
							"Something in the way...",
						),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.splitter.Split(tt.in)
			if out.Pre != tt.out.Pre {
				t.Fatalf("expected pre: %q, got %q", tt.out.Pre, out.Pre)
			}
			if len(out.Pages) != len(tt.out.Pages) {
				t.Fatalf("expected %d outputs, got %d", len(tt.out.Pages), len(out.Pages))
			}
			for i, o := range out.Pages {
				if o.Fm != tt.out.Pages[i].Fm {
					t.Fatalf("expected frontmatter: %q, got %q", tt.out.Pages[i].Fm, o.Fm)
				}
				if o.Md != tt.out.Pages[i].Md {
					t.Fatalf("expected markdown: %q, got %q", tt.out.Pages[i].Md, o.Md)
				}
			}
		})
	}
}
