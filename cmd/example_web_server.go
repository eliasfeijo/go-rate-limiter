package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	rlconfig "github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	Port     string       `mapstructure:"PORT"`
	LogLevel log.LogLevel `mapstructure:"LOG_LEVEL"`
}

var config = &Config{}

func main() {
	loadConfig()
	log.Log(log.Info, "Config loaded successfully")

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

	log.Log(log.Info, "Starting server on port "+config.Port)
	http.ListenAndServe(":"+config.Port, r)
}

func loadConfig() {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("LOG_LEVEL", log.Debug)

	err := viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			panic("Error reading config file")
		}
	}

	err = viper.Unmarshal(config)

	viper.UnmarshalKey("LOG_LEVEL", &config.LogLevel, viper.DecodeHook(logLevelHookFunc()))

	fmt.Println(config.LogLevel)

	log.SetLogger(NewLogger(config.LogLevel))

	err = rlconfig.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}
}

func logLevelHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Check that the data is string
		if f.Kind() == reflect.String {
			// Parse the string into a LogLevel
			logLevel, err := log.ParseLogLevel(data.(string))
			if err != nil {
				return nil, err
			}
			return logLevel, nil
		} else if f.Kind() == reflect.Uint8 {
			logLevel := log.LogLevel(data.(uint8))
			if logLevel < log.Debug {
				return log.Debug, nil
			} else if logLevel > log.Panic {
				return log.Panic, nil
			}
			return logLevel, nil
		}
		return log.Debug, fmt.Errorf("Invalid log level: %s", data)
	}
}
