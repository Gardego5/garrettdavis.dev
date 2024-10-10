//go:build static && fonts

package main

import (
	"embed"
	"net/http"

	"github.com/Gardego5/garrettdavis.dev/middleware"
)

//go:embed build/share/fonts
var fonts embed.FS

func init() {
	files.Handle("GET /fonts/", http.StripPrefix("/static/fonts/", middleware.FileServerFS(fonts)))
}
