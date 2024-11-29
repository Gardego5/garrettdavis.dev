package cookie

import "net/http"

var (
	State = http.Cookie{
		Name:     "oauth-state",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   300,
		SameSite: http.SameSiteNoneMode,
	}

	Session = http.Cookie{
		Name:     "garrettdavisdev-session",
		Path:     "/",
		HttpOnly: true,
		//Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
)
