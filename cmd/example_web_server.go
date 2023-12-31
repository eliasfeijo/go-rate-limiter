package main

import (
	"net/http"
	"os"

	rlconfig "github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

var port string

func main() {
	log.SetLogger(&Logger{})

	loadConfig()
	cfg := rlconfig.GetConfig()

	rateLimiterConfig := &limiter.RateLimiterConfig{
		MaxRequestsIpAddress:    cfg.RateLimiterIpAddressMaxRequests,
		LimitInSecondsIpAddress: cfg.RateLimiterIpAddressLimitInSeconds,
		BlockInSecondsIpAddress: cfg.RateLimiterIpAddressBlockInSeconds,
		TokensHeaderKey:         cfg.RateLimiterTokensHeaderKey,
		TokensConfig:            cfg.TokensConfig,
	}
	rateLimiterMiddleware := middleware.NewRateLimitMiddleware(rateLimiterConfig, cfg.RateLimiterStoreStrategy)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(rateLimiterMiddleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request accepted"))
	})

	log.Log(log.Info, "Starting server on port "+port)
	http.ListenAndServe(":"+port, r)
}

func loadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)

	err := viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			panic("Error reading config file")
		}
	}

	port = viper.GetString("PORT")

	err = rlconfig.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	log.Log(log.Info, "Config loaded successfully")
}
