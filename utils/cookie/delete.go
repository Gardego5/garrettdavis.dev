package cookie

import "net/http"

func Delete(w http.ResponseWriter, c http.Cookie) {
	c.MaxAge = -1
	c.Value = ""
	http.SetCookie(w, &c)
}
