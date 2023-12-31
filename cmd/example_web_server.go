package main

import (
	"net/http"

	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	loadConfig()

	rateLimiterConfig := &limiter.RateLimiterConfig{
		MaxRequestsIpAddress:    config.RateLimiterConfig.IpAddressMaxRequests,
		LimitInSecondsIpAddress: config.RateLimiterConfig.IpAddressLimitInSeconds,
		BlockInSecondsIpAddress: config.RateLimiterConfig.IpAddressBlockInSeconds,
		TokensHeaderKey:         config.RateLimiterConfig.TokensHeaderKey,
		MapTokenConfig:          config.RateLimiterConfig.MapTokenConfig,
	}
	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(rateLimiterConfig, config.RateLimiterConfig.StoreStrategy)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(rateLimiterMiddleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request accepted"))
	})

	log.Log(log.Info, "Starting server on port "+config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
