package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/messages"
	. "github.com/Gardego5/htmdsl"
	"github.com/elliotchance/pie/v2"
	"github.com/go-playground/validator/v10"
	"github.com/monoculum/formam"
)

type AdminMessages struct {
	messages *messages.Service
}

func NewAdminMessages(
	messages *messages.Service,
) *AdminMessages {
	return &AdminMessages{messages: messages}
}

func (h *AdminMessages) GET(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := access.Logger(ctx, "GetAdminUser")

	q := struct {
		Sort   messages.ListMessageInputSort `q:"sort"   validate:"oneof=ASC DESC"`
		Limit  int                           `q:"limit"  validate:"min=1,max=100"`
		Offset int                           `q:"offset" validate:"min=0"`
	}{messages.ListMessageInputSortDESC, 10, 0}
	r.ParseForm()
	access.Get[formam.Decoder](ctx).Decode(r.Form, &q)
	if err := access.Get[validator.Validate](ctx).Struct(q); err != nil {
		logger.Error("Error validating form", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Info("Listing messages", "query", q)

	msgs, err := h.messages.ListMessages(ctx, &messages.ListMessageInput{
		Sort: q.Sort, Limit: q.Limit, Offset: q.Offset})
	if err != nil {
		logger.Error("Error listing messages", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	count, err := h.messages.CountMessages(ctx)
	if err != nil {
		logger.Error("Error counting messages", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	list := Ul{Class("grid grid-cols-1 gap-6"),
		Attrs{
			"hx-target": "closest li",
			"hx-swap":   "outerHTML swap:0.1s",
		},
		pie.Map(msgs, func(msg model.ContactMessage) any {
			return Li{Class("relative rounded-sm border border-slate-500 bg-gray-800 p-4",
				"[&.htmx-swapping]:transition-opacity [&.htmx-swapping]:opacity-0 list-none",
			),
				Div{Class("flex justify-between mb-2 mx-2"),
					Span{Class("flex-grow"), msg.Name},
					A{Attrs{
						"href":   "mailto:" + msg.Email,
						"target": "_blank",
					}},
					msg.Email,
				},

				Pre{Class("bg-zinc-900 rounded-sm border border-slate-500 px-2 py-1 whitespace-pre-wrap"),
					msg.Message,
				},

				Div{Class("absolute -bottom-[7px] right-8 flex gap-2"),
					P{Class("rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-xs grid place-items-center"),
						msg.CreatedAt.Time.Format(time.RFC1123Z),
					},

					Button{
						Class("rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-red-800 grid place-items-center"),
						Attrs{"hx-delete": fmt.Sprintf("/admin/messages/%d", msg.ID)},
						Element("iconify-icon", Attrs{"icon": "mdi:delete-outline", "width": 20, "height": 20}),
					},
				},
			}
		}),
	}

	countDisplay := Div{Id("messages-count"),
		Class("inline-flex gap-2 justify-self-center col-span-full md:col-span-1 xs:row-start-2 md:row-start-1"),
		Button{Class("disabled:opacity-50"), Attrs{
			"disabled": AttrIf(q.Offset == 0),
			"name":     "offset",
			"value":    0,
		}, "<--"},

		Button{Class("disabled:opacity-50"), Attrs{
			"disabled": AttrIf(q.Offset == 0),
			"name":     "offset",
			"value":    max(q.Offset+q.Offset%q.Limit-q.Limit, 0),
		}, "<"},

		P{Class("text-gray-400"),
			q.Offset + 1, " - ", q.Offset + len(msgs), " of ", count,
		},

		Button{Class("disabled:opacity-50"), Attrs{
			"disabled": AttrIf(q.Offset+len(msgs) >= count),
			"name":     "offset",
			"value":    q.Offset + q.Offset%q.Limit + q.Limit,
		}, ">"},

		Button{Class("disabled:opacity-50"), Attrs{
			"disabled": AttrIf(q.Offset+len(msgs) >= count),
			"name":     "offset",
			"value":    count - count%q.Limit,
		}, "-->"},

		Input{"hidden": nil, "name": "offset", "value": q.Offset},
	}

	if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-Boosted") == "" {
		// Is this a rerender with new search parameters?
		w.Header().Set("HX-Replace-Url", r.URL.String())
		RenderContext(w, r.Context(), Fragment{
			list,
			append(countDisplay, Attrs{"hx-swap-oob": true}),
		})
	} else {
		// Or is this a fresh render after some kind of navigation?
		render.Page(w, r, nil, components.Header{Title: "Messages"}, components.Margins{
			Form{Class("grid grid-cols-2 md:grid-cols-3 items-center pb-4 px-4"),
				Attrs{
					"hx-params":  "*",
					"hx-swap":    "outerHTML swap:0.1s",
					"hx-target":  "next ul",
					"hx-trigger": "change",
					"hx-get":     "/admin/messages",
				},

				Div{Class("inline-block justify-self-start col-start-1 row-start-1"),
					Label{Attrs{"for": "messages-sort"}, "Sort: "},
					Select{Id("messages-sort"), Attrs{"name": "sort"},
						Option{Attrs{
							"value":    "DESC",
							"selected": AttrIf(q.Sort == messages.ListMessageInputSortASC),
						}, "Newest First"},
						Option{Attrs{
							"value":    "ASC",
							"selected": AttrIf(q.Sort == messages.ListMessageInputSortASC),
						}, "Oldest First"},
					},
				},

				countDisplay,

				Div{Class("inline-flex gap-4 align-items-center justify-self-end -col-start-2 row-start-1"),
					Label{Attrs{"for": "messages-limit"}, "Limit: "},
					Select{Id("messages-limit"), Attrs{"name": "limit"},
						Option{Attrs{"selected": AttrIf(q.Limit == 10)}, 10},
						Option{Attrs{"selected": AttrIf(q.Limit == 25)}, 25},
						Option{Attrs{"selected": AttrIf(q.Limit == 100)}, 100},
					},
				},
			},

			list,
		})
	}
}

func (h *AdminMessages) DELETE(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := access.Logger(ctx, "DeleteAdminMessage")

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		logger.Error("Error parsing id", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.messages.DeleteMessage(ctx, &messages.DeleteMessageInput{
		ID: int(id),
	}); err != nil {
		logger.Error("Error deleting message", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Message deleted", "id", id)
	w.WriteHeader(http.StatusOK)
}
