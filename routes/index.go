package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	. "github.com/Gardego5/htmdsl"
	"github.com/elliotchance/pie/v2"
)

func GetIndex(w http.ResponseWriter, r *http.Request) {
	middleware.RenderPage(r,
		Fragment{
			Title{"Garrett Davis"},
			Meta{{"name", "description"}, {"content", "Garrett Davis is a young software developer who cares deaply about creating great software for people."}},
		},

		components.Header{},
		components.Margins(
			Main{Class("text-lg"),
				P{Class("mb-12"),
					"hello! i'm garrett. ",
					"i'm learning a lot. ",
					"and i plan to keep at it. ",
				},

				Div{Class("my-12"),
					"i love making software that is: ",
					Ul{Class("inline md:inline-flex md:gap-5 items-center"),
						pie.Map([]struct{ adjective, class string }{
							{"enjoyable to use", "text-yellow-200"},
							{"reliable", "text-red-300"},
							{"maintainable", "text-blue-200"},
							{"simple", "text-green-300"},
						}, func(li struct{ adjective, class string }) any {
							return Fragment{" ", Li{Class(
								"relative after:absolute after:-right-3 after:top-1/2 after:-translate-y-1/2",
								"inline after:w-1 after:h-1 after:rounded-full after:bg-zinc-50",
								"after:content-none md:after:content-['_'] last:after:content-none",
								li.class,
							), li.adjective}}
						}),
						".",
					},
				},

				Div{Class("flex items-center justify-center gap-4 my-4"),
					Div{
						A{Attrs{
							{"target", "_blank"},
							{"rel", "noopener noreferrer"},
							{"href", "https://github.com/Gardego5"},
							{"aria-label", "View Garrett's GitHub Profile"},
							{"class", "border border-slate-500 bg-zinc-950 border-dotted p-2 pl-2 sm:pl-0.5 flex text-baseline h-8 items-center text-sm gap-1 text-nowrap"}},
							Element("iconify-icon", Attrs{
								{"icon", "mdi:github"},
								{"width", "32"}, {"height", "32"},
								{"class", "sm:scale-75"},
							}),
							Span{Class("hidden sm:inline"), "github"},
						},
					},

					Span{Class("text-sm text-center"),
						Span{Class("text-nowrap"), "&lt;--"},
						" find me here ",
						Span{Class("text-nowrap"), "--&gt;"},
					},

					Div{
						A{Attrs{
							{"target", "_blank"},
							{"rel", "noopener noreferrer"},
							{"href", "https://www.linkedin.com/in/garrett-davis-8793a721b/"},
							{"aria-label", "View Garrett's LinkedIn Profile"},
							{"class", "border border-slate-500 bg-zinc-950 border-dotted p-2 pr-2 sm:pr-0.5 flex text-baseline h-8 items-center text-sm gap-1 text-nowrap"}},
							Span{Class("hidden sm:inline"), "linkedin"},
							Element("iconify-icon", Attrs{
								{"icon", "mdi:linkedin"},
								{"width", "32"}, {"height", "32"},
								{"class", "sm:scale-75"},
							}),
						},
					},
				},

				Hr{},

				pie.Map(livePostsList(), func(post BlogPost) RenderedHTML {
					return Fragment{post.Summary(), Hr{}}
				}),
			},
		),
	)
}
