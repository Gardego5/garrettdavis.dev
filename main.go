package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Gardego5/garrettdavis.dev/resource/initialize"
	"github.com/Gardego5/garrettdavis.dev/resource/middleware"
	"github.com/Gardego5/garrettdavis.dev/routes"
	"github.com/Gardego5/garrettdavis.dev/service/blog"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	"github.com/Gardego5/garrettdavis.dev/service/messages"
	"github.com/Gardego5/garrettdavis.dev/service/object"
	"github.com/Gardego5/garrettdavis.dev/service/presentations"
	"github.com/Gardego5/garrettdavis.dev/utils"
	"github.com/Gardego5/garrettdavis.dev/utils/mux"
	"github.com/Gardego5/garrettdavis.dev/utils/symetric"
	"github.com/Gardego5/goutils/env"
	"github.com/go-playground/validator/v10"
)

var (
	Env = utils.Must(env.Load[struct {
		ApplicationSecret string        `env:"APPLICATION_SECRET" validate:"required"`
		GithubOauthId     string        `env:"GITHUB_OAUTH_CLIENT_ID" validate:"required"`
		GithubOauthSecret string        `env:"GITHUB_OAUTH_CLIENT_SECRET" validate:"required"`
		Host              string        `env:"HOST=0.0.0.0" validate:"required"`
		LogLevel          slog.LevelVar `env:"LOG_LEVEL=INFO"`
		Port              int           `env:"PORT=8080" validate:"required"`
		RedisUrl          string        `env:"REDIS_URL" validate:"required"`
		TursoAuthToken    string        `env:"TURSO_AUTH_TOKEN" validate:"required"`
		TursoDatabaseUrl  string        `env:"TURSO_DATABASE_URL" validate:"required"`
		ImagesBucket      string        `env:"IMAGES_BUCKET" validate:"required"`
	}]())

	Validate = validator.New()
	Logger   = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: &Env.LogLevel}))
	DB       = initialize.NewDB(Env.TursoDatabaseUrl, Env.TursoAuthToken)
	Redis    = initialize.NewRedis(Env.RedisUrl)
	Enforcer = utils.Must(initialize.Enforcer(DB))
	Caches   = initialize.Caches(Redis)
	Block    = symetric.Block(Env.ApplicationSecret)

	// services
	Blog          = blog.New()
	CurrentUser   = currentuser.New(Caches)
	Messages      = utils.Must(messages.New(DB))
	Presentations = presentations.New()
	ImagesBucket  = utils.Must(object.New(context.Background(), Env.ImagesBucket, Logger))

	//go:embed static
	Static embed.FS
	//go:embed build/share/fonts
	Fonts embed.FS

	Mux = mux.NewServeMux(func(m *mux.ServeMux) {
		m.Group("/admin", func(m *mux.ServeMux) {
			m.Group("/messages", func(m *mux.ServeMux) {
				h := routes.NewAdminMessages(Messages)
				m.HandleFunc("GET", h.GetAdminMessage)
				m.HandleFunc("DELETE /{id}", h.DeleteAdminMessage)
			})
			m.Handle("GET /user", routes.NewAdminUser(CurrentUser))
			m.Group("/coffee", func(m *mux.ServeMux) {
				h := routes.NewAdminCoffee(ImagesBucket)
				m.HandleFunc("GET", h.GetAdminCoffee)
			})
		})

		m.Group("/auth", func(m *mux.ServeMux) {
			m.Handle("GET /callback", routes.NewAuthCallback(
				Env.GithubOauthId, Env.GithubOauthSecret,
				Validate, Block, CurrentUser, Caches))
			m.Group("/signin", func(m *mux.ServeMux) {
				h := routes.NewAuthSignin(Env.GithubOauthId, Block)
				m.HandleFunc("GET", h.GET)
				m.HandleFunc("POST", h.POST)
			})
			m.Handle("POST /signout", routes.NewAuthSignout())
		})

		m.Handle("GET /blog/{slug}", routes.NewBlog(Blog))

		m.Group("/contact", func(m *mux.ServeMux) {
			h := routes.NewContact(Messages, Validate)
			m.HandleFunc("GET", h.GET)
			m.HandleFunc("POST", h.POST)
		})

		m.Handle("GET /presentations/{slug}", routes.NewPresentations(Presentations))

		m.Handle("GET /resume", routes.NewResume())
		m.Group("/static", func(m *mux.ServeMux) {
			m.Handle("GET /", middleware.FileServerFS(Static))
			m.Handle("GET /fonts/", http.StripPrefix("/static/fonts/", middleware.FileServerFS(Fonts)))
		},
			middleware.NeuteredFileSystem,
		)
		m.Handle("GET /{$}", routes.NewIndex(Blog))
		m.HandleFunc("GET /", routes.Get404)
	},
		middleware.LoggerAndSessions(Logger, true),
		middleware.TrailingSlash,
		middleware.GenericAssets4(
			Blog, CurrentUser, Messages, Presentations,
		),
		middleware.GenericAsset(
			Validate,
		),
	)
)

func main() {
	defer func() {
		if err := errors.Join(
			DB.Close(),
			Messages.Close(),
		); err != nil {
			Logger.Error("Error cleaning up resources", "error", err)
			os.Exit(1)
		}
	}()

	addr := fmt.Sprintf("%s:%d", Env.Host, Env.Port)
	if err := http.ListenAndServe(addr, Mux); err != nil {
		Logger.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
