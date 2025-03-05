package routes

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/messages"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
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

func (h *AdminMessages) GetAdminMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := access.Logger(ctx, "GetAdminUser")

	q := struct {
		Sort   messages.ListMessageInputSort `q:"sort"   validate:"oneof=ASC DESC"`
		Limit  int                           `q:"limit"  validate:"min=1,max=100"`
		Offset int                           `q:"offset" validate:"min=0"`
	}{"ASC", 10, 0}
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

	var isFirstRender bool
	{
		currentUrlHeader := r.Header.Get("HX-Current-URL")
		url, err := url.Parse(currentUrlHeader)
		switch {
		case currentUrlHeader == "":
			isFirstRender = true
		case err != nil:
			logger.Warn("Error parsing current URL", "error", err)
			isFirstRender = true
		default:
			isFirstRender = url.Path != r.URL.Path
		}
	}

	countDisplay := Div{Id("messages-count"), Class("text-gray-400"),
		"Showing ",
		Span{Id("messages-count.range-start"), q.Offset + 1},
		"-",
		Span{Id("messages-count.range-end"), q.Offset + len(msgs)},
		" of ",
		Span{Id("messages-count.total"), count},
	}

	if isFirstRender {

		render.Page(w, r, nil, components.Header{Title: "Messages"}, components.Margins{
			Form{Class("flex justify-between items-center pb-4"),
				Attrs{
					"hx-params":  "*",
					"hx-swap":    "outerHTML swap:0.1s",
					"hx-target":  "next ul",
					"hx-trigger": "submit, change delay:500ms",
				},

				Div{Class("inline-block"),
					Label{Attrs{"for": "messages-sort"}, "Sort: "},
					Select{Id("messages-sort"), Attrs{
						"name":   "sort",
						"hx-get": "/admin/messages",
					},
						Option{If(q.Sort == messages.ListMessageInputSortASC, Attrs{"selected": nil}).Default(),
							Attrs{"value": "ASC"}, "Oldest First"},
						Option{If(q.Sort == messages.ListMessageInputSortDESC, Attrs{"selected": nil}).Default(),
							Attrs{"value": "DESC"}, "Newest First"},
					},
				},

				Div{Class("inline-block"),
					Label{Attrs{"for": "messages-limit"}, "Limit: "},
					Select{Id("messages-limit"), Attrs{
						"name":    "limit",
						"hx-get":  "/admin/messages",
						"x-data":  fmt.Sprintf("{ selected: %d }", q.Limit),
						"x-model": "selected",
					},
						Option{Attrs{"value": "10", ":selected": "selected === 10"}, "10"},
						Option{Attrs{"value": "25", ":selected": "selected === 25"}, "25"},
						Option{Attrs{"value": "100", ":selected": "selected === 100"}, "100"},
					},
				},

				countDisplay,
			},

			list,
		})
	} else {
		w.Header().Set("HX-Push-Url", r.URL.String())
		RenderContext(w, r.Context(), Fragment{
			list,
			append(countDisplay, Attrs{"hx-swap-oob": true}),
		})
	}
}

func (h *AdminMessages) DeleteAdminMessage(w http.ResponseWriter, r *http.Request) {
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
