package routes

import (
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Gardego5/garrettdavis.dev/components"
	"github.com/Gardego5/garrettdavis.dev/model"
	"github.com/Gardego5/garrettdavis.dev/resource/access"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	"github.com/Gardego5/garrettdavis.dev/utils/bimarshal"
	"github.com/Gardego5/garrettdavis.dev/utils/cookie"
	"github.com/Gardego5/garrettdavis.dev/utils/symetric"
	. "github.com/Gardego5/htmdsl"
	"github.com/go-playground/validator/v10"
)

type AuthCallback struct {
	clientId, clientSecret string
	validate               *validator.Validate
	block                  cipher.Block
	currentuser            *currentuser.Service
	subjects               bimarshal.Cache[model.Subject]
	accessTokens           bimarshal.Cache[model.GHAccessToken]
}

func NewAuthCallback(
	clientId, clientSecret string,
	validator *validator.Validate,
	block cipher.Block,
	currentuser *currentuser.Service,
	caches bimarshal.RegisteredCaches,
) *AuthCallback {
	return &AuthCallback{
		clientId:     clientId,
		clientSecret: clientSecret,
		validate:     validator,
		block:        block,
		currentuser:  currentuser,
		subjects:     bimarshal.Get[model.Subject](caches),
		accessTokens: bimarshal.Get[model.GHAccessToken](caches),
	}
}

func (h *AuthCallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := access.Logger(ctx, "PostAuthCallback")

	// delete the cookie regardless of the outcome
	cookie.Delete(w, cookie.State)

	// read the state cookie
	c, err := r.Cookie(cookie.State.Name)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn("Error reading state cookie", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	// decrypt the state
	state, err := symetric.Decrypt(h.block, c.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn("Error decrypting state cookie", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	// parse the query parameters we are given from github
	q := r.URL.Query()
	payload := struct {
		Code  string `validate:"required"`
		State string `validate:"required"`
	}{q.Get("code"), q.Get("state")}
	if err = h.validate.Struct(payload); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn("Error validating payload", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	// compare the state we received from github to the state we sent to github
	if state != payload.State {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn("State mismatch", "expected", state, "received", payload.State)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	// send a request to github to get the access token
	q = url.Values{}
	q.Set("client_id", h.clientId)
	q.Set("client_secret", h.clientSecret)
	q.Set("code", payload.Code)
	q.Set("redirect_uri", "http://localhost:8080/auth/callback")
	u := fmt.Sprintf("https://github.com/login/oauth/access_token?%s", q.Encode())
	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error creating request to github", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error sending request to github", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	oauth := model.GHAccessToken{}

	if err = json.NewDecoder(res.Body).Decode(&oauth); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error decoding github response", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	if err = h.validate.Struct(oauth); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error validating github response", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	user, err := h.currentuser.GetUserByAccessToken(ctx, oauth.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error getting user from github", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	sub := model.Subject{User: user.GetLogin()}

	/*
		TODO: signup restrictions.
		if has, err := access.Enforcer(ctx).
			Enforce(sub, access.Path(ctx), r.Method); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("Error enforcing policy", "error", err)
			render.Page(w, r, nil, P{"An error has occurred."})
			return
		} else if !has {
			w.WriteHeader(http.StatusForbidden)
			logger.Warn("User does not have permission")
			render.Page(w, r, nil, P{"Why are you here?"})
			return
		}
	*/

	sessionCookie := cookie.Session
	sessionCookie.Value = access.Session(ctx)
	sessionCookie.MaxAge = 60 * 60 * 5 // 5 hours
	logger.Info("Creating session", "session", sessionCookie.Value, "oauth", oauth)

	cacheDuration := time.Duration(sessionCookie.MaxAge) * time.Second

	if err = h.accessTokens.Set(ctx, sessionCookie.Value, oauth, cacheDuration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error storing session in cache", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	if err = h.subjects.Set(ctx, sessionCookie.Value, sub, cacheDuration); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Error("Error storing subject in cache", "error", err)
		render.Page(w, r, nil, P{"An error has occurred."})
		return
	}

	http.SetCookie(w, &sessionCookie)
	w.WriteHeader(http.StatusOK)
	render.Page(w, r, nil,
		components.Header{},
		components.Margins(P{"Authenticated!"}))
}
