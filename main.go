package main

import (
	"fmt"
	"net/http"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	cfg := config.GetConfig()

	rateLimiterConfig := &limiter.RateLimiterConfig{
		MaxRequestsIpAddress:    cfg.RateLimiterIpAddressMaxRequests,
		LimitInSecondsIpAddress: cfg.RateLimiterIpAddressLimitInSeconds,
		BlockInSecondsIpAddress: cfg.RateLimiterIpAddressBlockInSeconds,
		TokensHeaderKey:         cfg.RateLimiterTokensHeaderKey,
		TokensConfig:            cfg.TokensConfig,
	}
	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(rateLimiterConfig)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(rateLimiterMiddleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request accepted!")
	})

	http.ListenAndServe(":"+cfg.Port, r)
}
