package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Gardego5/garrettdavis.dev/global"
	"github.com/Gardego5/garrettdavis.dev/middleware"
	"github.com/Gardego5/garrettdavis.dev/routes"
	"github.com/Gardego5/goutils/env"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var (
	db     = middleware.NewDB()
	mux    = middleware.NewServeMux(middleware.Logger, middleware.LogRequests)
	common = mux.Use(middleware.Markdown, middleware.WriterRef, middleware.DB(db))
)

func init() {
	common.HandleFunc("GET /", routes.Get404)
	common.HandleFunc("GET /{$}", routes.GetIndex)
	common.HandleFunc("GET /blog/{slug}", routes.GetBlogPage)
	common.HandleFunc("GET /contact", routes.GetContact)
	common.HandleFunc("GET /presentation/{slug}", routes.GetPresentation)
	common.HandleFunc("GET /resume", routes.GetResume)
	common.HandleFunc("POST /contact", routes.PostContact)
}

func main() {
	env := env.MustLoad[struct {
		Port string `env:"PORT=8080"`
		Host string `env:"HOST=0.0.0.0"`
	}]()

	defer func() {
		if err := db.Close(); err != nil {
			global.Logger.Error("Error closing database connection", "error", err)
			os.Exit(1)
		}
	}()

	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		lambda.Start(httpadapter.NewFunctionURL(mux).ProxyWithContext)
	} else {
		http.ListenAndServe(fmt.Sprintf("%s:%s", env.Host, env.Port), mux)
	}
}
