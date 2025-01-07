package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/messages"
	. "github.com/Gardego5/htmdsl"
	"github.com/elliotchance/pie/v2"
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

	messages, err := h.messages.ListMessages(ctx, &messages.ListMessageInput{
		Limit: 10, Offset: 0,
	})
	if err != nil {
		logger.Error("Error listing messages", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.Page(w, r, nil, components.Header{Title: "Messages"}, components.Margins(
		Ul{Attrs{
			Class("grid grid-cols-1 gap-6"),
			{"hx-confirm", "Are you sure you want to delete this message?"},
			{"hx-target", "closest li"},
			{"hx-swap", "outerHTML swap:1s"},
		},
			pie.Map(messages, func(msg model.ContactMessage) any {
				return Li{Class(
					"relative rounded-sm border border-slate-500 bg-gray-800 p-4",
					"[&.htmx-swapping]:transition-opacity [&.htmx-swapping]:opacity-0",
				),
					Div{Class("flex justify-between mb-2 mx-2"),
						Span{msg.Name},
						A{Attrs{{"href", "mailto:" + msg.Email}, {"target", "_blank"}},
							msg.Email,
						},
					},

					P{Class("bg-zinc-900 rounded-sm border border-slate-500 px-2 py-1"),
						msg.Message,
					},

					Div{Class("absolute -bottom-[7px] right-8 flex gap-2"),
						P{Class("rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-xs grid place-items-center"),
							msg.CreatedAt.Time.Format("Mon, July 30, 2006 15:04:05"),
						},

						Button{Attrs{
							Class("rounded-sm border border-slate-500 bg-zinc-900 px-4 py-1 text-sm hover:bg-red-800 grid place-items-center"),
							{"hx-delete", fmt.Sprintf("/admin/messages/%d", msg.ID)},
						},
							Element("iconify-icon", Attrs{{"icon", "mdi:delete-outline"}, {"width", "20"}, {"height", "20"}}),
						},
					},
				}
			}),
		},
	))
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
