package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Gardego5/garrettdavis.dev/resource/initialize"
	"github.com/Gardego5/garrettdavis.dev/resource/middleware"
	"github.com/Gardego5/garrettdavis.dev/resource/render"
	"github.com/Gardego5/garrettdavis.dev/routes"
	"github.com/Gardego5/garrettdavis.dev/service/blog"
	"github.com/Gardego5/garrettdavis.dev/service/currentuser"
	"github.com/Gardego5/garrettdavis.dev/service/messages"
	"github.com/Gardego5/garrettdavis.dev/service/object"
	"github.com/Gardego5/garrettdavis.dev/service/presentations"
	"github.com/Gardego5/garrettdavis.dev/service/resume"
	"github.com/Gardego5/garrettdavis.dev/utils"
	"github.com/Gardego5/garrettdavis.dev/utils/mux"
	"github.com/Gardego5/garrettdavis.dev/utils/symetric"
	"github.com/Gardego5/goutils/env"
	"github.com/go-playground/validator/v10"
	"github.com/monoculum/formam"
)

var (
	Env = utils.Must(env.Load[struct {
		ApplicationSecret string        `env:"APPLICATION_SECRET" validate:"required"`
		BaseUrl           string        `env:"BASE_URL=https://garrettdavis.dev" validate:"required"`
		GithubOauthId     string        `env:"GITHUB_OAUTH_CLIENT_ID" validate:"required"`
		GithubOauthSecret string        `env:"GITHUB_OAUTH_CLIENT_SECRET" validate:"required"`
		Host              string        `env:"HOST=0.0.0.0" validate:"required"`
		ImagesBucket      string        `env:"IMAGES_BUCKET" validate:"required"`
		LogLevel          slog.LevelVar `env:"LOG_LEVEL=INFO"`
		Port              int           `env:"PORT=8080" validate:"required"`
		RedisUrl          string        `env:"REDIS_URL" validate:"required"`
		TursoAuthToken    string        `env:"TURSO_AUTH_TOKEN" validate:"required"`
		TursoDatabaseUrl  string        `env:"TURSO_DATABASE_URL" validate:"required"`
	}]())

	Validate = validator.New()
	Logger   = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &Env.LogLevel}))
	DB       = initialize.NewDB(Env.TursoDatabaseUrl, Env.TursoAuthToken)
	Redis    = initialize.NewRedis(Env.RedisUrl)
	Enforcer = utils.Must(initialize.Enforcer(DB, Logger))
	Caches   = initialize.Caches(Redis)
	Block    = symetric.Block(Env.ApplicationSecret)

	// services
	Blog          = blog.New()
	CurrentUser   = currentuser.New(Caches)
	ImagesBucket  = utils.Must(object.New(context.Background(), Env.ImagesBucket, Logger))
	Messages      = utils.Must(messages.New(DB))
	Presentations = presentations.New()
	Resume        = resume.New(Validate)

	// these are assets / configuration included at build time
	//go:embed build
	Build        embed.FS
	CacheID      string
	Static       = utils.Must(fs.Sub(Build, "build"))
	StaticPrefix = fmt.Sprintf("/static/%s", CacheID)

	Mux = mux.NewServeMux(func(m *mux.ServeMux) {
		m.Group("/admin", func(m *mux.ServeMux) {
			m.Group("/messages", func(m *mux.ServeMux) {
				h := routes.NewAdminMessages(Messages)
				m.HandleFunc("GET", h.GET)
				m.HandleFunc("DELETE /{id}", h.DELETE)
			})
			m.Handle("GET /user", routes.NewAdminUser(CurrentUser))
			m.Group("/coffee", func(m *mux.ServeMux) {
				h := routes.NewAdminCoffee(ImagesBucket)
				m.HandleFunc("GET", h.GetAdminCoffee)
			})
		},
			middleware.Authorization(Logger, Enforcer, CurrentUser))

		m.Group("/auth", func(m *mux.ServeMux) {
			m.Handle("GET /callback", routes.NewAuthCallback(
				Env.GithubOauthId, Env.GithubOauthSecret,
				Validate, Block, CurrentUser, Caches, Enforcer, Env.BaseUrl))
			m.Group("/signin", func(m *mux.ServeMux) {
				h := routes.NewAuthSignin(Env.GithubOauthId, Block, Env.BaseUrl)
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

		m.Handle("GET "+StaticPrefix+"/", http.StripPrefix(StaticPrefix, middleware.FileServerFS(Static)),
			middleware.FileSystem,
			mux.MiddlewareFunc(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Cache-Control", "max-age=31536000, immutable")
					next.ServeHTTP(w, r)
				})
			}),
		)

		m.Handle("GET /{$}", routes.NewIndex(Blog))
		m.HandleFunc("GET /", routes.Get404)
	},
		middleware.LoggerAndSessions(Logger, true),
		middleware.TrailingSlash,
		middleware.Inject(
			middleware.Syringe(Blog),
			middleware.Syringe(CurrentUser),
			middleware.Syringe(Messages),
			middleware.Syringe(Presentations),
			middleware.Syringe(Resume),
			middleware.Syringe(Validate),
			middleware.Syringe(utils.Ptr(render.StaticPathPrefix(StaticPrefix))),
			middleware.Syringe(formam.NewDecoder(&formam.DecoderOptions{TagName: "q"})),
		),
	)
)

func main() {
	Logger := Logger.With("function", "main")

	defer func() {
		if err := errors.Join(
			DB.Close(),
			Messages.Close(),
		); err != nil {
			Logger.Error("Error cleaning up resources", "error", err)
			os.Exit(1)
		}
	}()

	tickPolicy := time.NewTicker(time.Minute)

	chServer := make(chan struct{})
	go func() {
		addr := fmt.Sprintf("%s:%d", Env.Host, Env.Port)
		if err := http.ListenAndServe(addr, Mux); err != nil {
			Logger.Error("Error starting server", "error", err)
		}
		chServer <- struct{}{}
	}()

	for {
		select {
		case <-tickPolicy.C:
			if err := Enforcer.LoadPolicy(); err != nil {
				Logger.Error("Error loading policy", "error", err)
				return
			} else {
				Logger.Info("Policy reloaded")
			}

		case <-chServer:
			return
		}
	}
}
