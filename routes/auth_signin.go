package routes

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/utils/cookie"
	"github.com/Gardego5/garrettdavis.dev/utils/symetric"
	. "github.com/Gardego5/htmdsl"
	. "github.com/Gardego5/htmdsl/util"
)

type AuthSignin struct {
	clientId string
	block    cipher.Block
	baseUrl  string
}

func NewAuthSignin(
	clientId string,
	block cipher.Block,
	baseUrl string,
) *AuthSignin {
	return &AuthSignin{
		clientId: clientId,
		block:    block,
		baseUrl:  baseUrl,
	}
}

func (h *AuthSignin) GET(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(cookie.State.Name)

	render.Page(w, r, nil,
		components.Header{},
		components.Margins(Form{Attrs{{"method", "post"}, {"action", "/auth/signin"}},
			Button{Attrs{{"type", "submit"}}, "Sign in"},
			If(c != nil, func() any {
				value, _ := symetric.Decrypt(h.block, c.Value)
				return P{"Cookie: ", base64.RawURLEncoding.EncodeToString([]byte(value))}
			}),
		}))
}

func (h *AuthSignin) POST(w http.ResponseWriter, r *http.Request) {
	logger := access.Logger(r.Context(), "PostAuth")

	codeBytes := make([]byte, 40)
	if _, err := rand.Read(codeBytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error generating code for authentication", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	state := base64.RawURLEncoding.EncodeToString(codeBytes)

	encryptedCode, err := symetric.Encrypt(h.block, state)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error encrypting code for authentication", "error", err)
		render.Page(w, r, nil,
			components.Header{},
			components.Margins(P{"An error has occurred."}))
		return
	}

	c := cookie.State
	c.Value = encryptedCode
	http.SetCookie(w, &c)

	redirectUri := h.baseUrl + "/auth/callback"
	logger.Info("Redirecting to github for authentication", "redirect_uri", redirectUri)
	q := url.Values{}
	q.Set("client_id", h.clientId)
	q.Set("redirect_uri", redirectUri)
	q.Set("allow_signup", "false")
	q.Set("state", state)
	u := fmt.Sprintf("https://github.com/login/oauth/authorize?%s", q.Encode())
	render.Page(w, r, nil,
		components.Header{},
		components.Margins{Attrs{{"x-data"}, {"@click", "$refs.authorize.click()"}},
			P{"Redirecting to github for authentication..."},
			P{"If you are not redirected, click the link below."},
			A{Attrs{{"href", u}, {"x-ref", "authorize"}}, "Sign in with Github"},
		})
}
