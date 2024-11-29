package components

import (
	"context"

	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	"github.com/Gardego5/htmdsl"
	"github.com/Gardego5/htmdsl/util"
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

func (header Header) Render(ctx context.Context) html.RenderedHTML {
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
			func(ctx context.Context) any {
				svc := access.Get[currentuser.Service](ctx)
				session := access.Session(ctx)
				if user, err := svc.GetUserBySession(ctx, session); err == nil {
					identifier := user.GetLogin()
					if identifier == "" {
						identifier = "user"
					}
					return html.Fragment{
						html.Form{html.Attrs{{"method", "POST"}, {"action", "/auth/signout"}},
							html.Button{"signout"},
						},

						html.Div{html.Attrs{html.Class("mr-auto relative"),
							{"x-data", `{
open:false,
toggle() {
	if (this.open) return this.close()
	this.$refs.button.focus()
	this.open = true
},
close(focusAfter) {
	if (! this.open) return
	this.open = false
	focusAfter && focusAfter.focus()
},
get menuStyle() {
	var rect = this.$refs.button.getBoundingClientRect()
	return { "--x-pos": rect.left+"px", "--y-pos": rect.bottom+"px" }
}
}`},
							{"@keydown.escape.prevent.stop", "close($refs.button)"},
							{"@focusin.window", "! $refs.panel.contains($event.target) && close()"},
							{"x-id", "['dropdown-button']"},
						},
							html.Button{html.Attrs{
								{"x-ref", "button"},
								{"@click", "toggle()"},
								{":aria-expanded", "open"},
								{":aria-controls", "$id('dropdown-button')"},
							},
								"menu"},

							html.Template{html.Attrs{
								{"x-teleport", "body"},
							},
								html.Menu{html.Attrs{
									{"x-ref", "panel"},
									{"x-show", "open"},
									{"x-transition.origin.top.left"},
									{":id", "$id('dropdown-button')"},
									{"style", "display: none;"},
									{":style", "menuStyle"},
									{"x-anchor.bottom-start.offset.10", "$refs.button"},
									{"@click.outside", "close($refs.button)"},
									html.Class("bg-zinc-950 shadow-md text-blue z-20 p-2"),
								},
									html.Li{html.A{html.Attrs{{"href", "/admin/user"}}, identifier}},
									html.Li{html.A{html.Attrs{{"href", "/admin/messages"}}, "messages"}},
								},
							},
						},
					}
				} else {
					return html.Form{html.Attrs{html.Class("mr-auto"), {"method", "POST"}, {"action", "/auth/signin"}},
						html.Button{"signin"},
					}
				}
			},

			html.A{html.Attrs{{"href", "/contact"}}, "contact"},
			html.A{html.Attrs{{"href", "/resume"}}, "resume"},
		},
	}.Render(ctx)
}

type Margins html.Fragment

func (children Margins) Render(ctx context.Context) html.RenderedHTML {
	return html.Div{html.Class("p-4 m-auto max-w-3xl"), []any(children)}.Render(ctx)
}
