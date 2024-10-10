package components

import (
	"github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/util"
	"github.com/elliotchance/pie/v2"
)

type (
	headerType int
	Header     struct {
		Variant headerType
		Title   string
	}
)

var _ html.HTML = (*Header)(nil)

const (
	HeaderFloating headerType = iota
	HeaderStatic
)

func (header Header) Render() html.RenderedHTML {
	headerTextClass := "mx-auto max-w-3xl px-4 pb-1 md:py-1"
	return html.Header{
		html.Class("overflow-hidden border-b border-slate-500 bg-zinc-950 pt-3 top-0 z-10"),
		util.If(header.Variant == HeaderFloating, html.Class("sticky")).Else(html.Class("relative")),

		util.If(header.Title != "",
			html.Div{html.Class(headerTextClass, "relative flex items-start justify-between gap-4"),
				html.H1{html.Class("text-xl md:text-3xl"), header.Title},
				html.A{html.Attrs{
					{"class", "md:absolute md:left-[calc(100%-1rem)] pr-2 lg:left-2 md:top-0 md:-translate-x-full text-right text-sm lg:text-base"},
					{"href", "/"}},
					"Garrett", html.Br{}, "Davis",
				},
			}).Else(
			html.Div{html.Class(headerTextClass, "text-xl md:text-3xl"),
				html.A{html.Attrs{{"href", "/"}}, "Garrett Davis"},
			}),

		html.Nav{
			html.Class("relative -left-1 mx-auto mb-1 flex justify-end max-w-3xl flex-shrink items-center gap-4 px-3 text-sm after:absolute after:-left-[100%] after:-z-10 after:h-[1px] after:w-[1000vw] after:bg-slate-500 md:text-base [&>*]:bg-zinc-950 [&>*]:px-2"),
			pie.Map([]struct{ href, name string }{
				//{"/blog", "blog"},
				{"/contact", "contact"},
				{"/resume", "resume"},
			}, func(val struct{ href, name string }) any {
				return html.A{html.Attrs{{"href", val.href}}, val.name}
			}),
		},
	}.Render()
}

func Margins(children ...any) html.RenderedHTML {
	return html.Div{html.Class("p-4 m-auto max-w-3xl"), children}.Render()
}
