package routes

import (
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/utils/cookie"
)

type AuthSignout struct{}

func NewAuthSignout() *AuthSignout {
	return &AuthSignout{}
}

func (*AuthSignout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// session := middleware.GetSession(r.Context())
	cookie.Delete(w, cookie.Session)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
