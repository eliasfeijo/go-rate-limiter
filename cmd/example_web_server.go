package main

import (
	"net/http"

	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	loadConfig()

	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(&config.RateLimiterConfig)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(rateLimiterMiddleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request accepted"))
	})

	log.Log(log.Info, "Starting server on port "+config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
