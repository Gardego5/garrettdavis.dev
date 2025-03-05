package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/object"
	. "github.com/Gardego5/htmdsl"
)

type AdminCoffee struct {
	images *object.Service
}

func NewAdminCoffee(
	images *object.Service,
) *AdminCoffee {
	return &AdminCoffee{images: images}
}

func (h *AdminCoffee) GetAdminCoffee(w http.ResponseWriter, r *http.Request) {
	render.Page(w, r, Title{"Admin Coffee"},
		components.Header{Title: "Admin Coffee"},
		components.Margins(Form{Class("relative grid gap-2 rounded-sm border border-slate-500 bg-gray-800 p-4 sm:grid-cols-2 md:grid-cols-3"),
			Attrs{"x-data": `{images:[]}`},
			H1{"Create a new Coffee Post"},

			Div{Class("col-span-full"),
				Label{Class("pl-2 text-slate-300 text-sm italic"), "Title"},
				Input{
					"class":       "block w-full px-2 py-1 rounded-sm border bg-zinc-900 border-slate-500",
					"name":        "title",
					"required":    nil,
					"type":        "text",
					"placeholder": "Sumatra just ain't it",
				},
				P{Id("title-error"), Class("text-sm text-red-600 px-2 empty:hidden")},
			},

			Div{Class("col-span-full grid place-items-center"),
				Button{Class("rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-slate-800"),
					Attrs{"@click": `$data.images.push({id: window.crypto.randomUUID()})`},
					"Add Image",
				},
			},

			Template{
				Attrs{"x-for": "image in images", ":key": "image.id"},
				Div{Class("col-span-full grid"),
					Div{Class("col-span-full"),
						Label{Class("pl-2 text-slate-300 text-sm italic"), "Picture Description"},
						Textarea{Class("block w-full px-2 py-1 rounded-sm border bg-zinc-900 border-slate-500"),
							Attrs{"name": "message", "required": nil, "placeholder": "Picture "},
						},
						P{Id("message-error"), Class("text-sm text-red-600 px-2 empty:hidden")},
					},

					Div{Class("col-span-full"),
						Input{
							"class":    "file:px-2 file:py-1 file:rounded-sm file:border file:border-solid file:bg-zinc-900 file:border-slate-500",
							"type":     "file",
							":name":    `"image-" + image.id`,
							"accept":   "image/*",
							"required": nil,
						},
					},
				},
			},

			Button{Class("absolute -bottom-[7px] right-8 rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-slate-800"),
				Attrs{"type": "submit"},
				"Send",
			},
		}))
}
