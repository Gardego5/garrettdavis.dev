package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	. "github.com/Gardego5/htmdsl"
)

type AdminUser struct {
	currentuser *currentuser.Service
}

func NewAdminUser(
	currentuser *currentuser.Service,
) *AdminUser {
	return &AdminUser{currentuser: currentuser}
}

func (h *AdminUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := access.Logger(ctx, "GetAdminUser")
	session := access.Session(ctx)

	user, err := h.currentuser.GetUserBySession(ctx, session)
	if err != nil {
		logger.Error("Error getting user", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.Page(w, r, nil,
			components.Header{},
			components.Margins(P{"An error has occurred."}))
	}

	data, _ := json.MarshalIndent(user, "", "  ")
	render.Page(w, r, nil,
		components.Header{Title: "User"},
		components.Margins(
			Pre{Code{
				string(data),
			}},
		))
}
