//go:build static

//go:generate tailwindcss -c ./tailwind.config.js -i ./input.css -o ./static/css/style.css --minify

package main

import (
	"embed"

	"github.com/Gardego5/garrettdavis.dev/middleware"
)

var files = mux.Group("/static", middleware.NeuteredFileSystem)

//go:embed static
var static embed.FS

func init() {
	files.Handle("GET /", middleware.FileServerFS(static))
}
