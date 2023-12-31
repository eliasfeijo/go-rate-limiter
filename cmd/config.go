package main

import (
	"fmt"
	"os"
	"reflect"

	rlconfig "github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	// Server port
	Port string `mapstructure:"PORT"`
	// Log level (debug, info, warn, error, panic)
	LogLevel log.LogLevel `mapstructure:"LOG_LEVEL"`

	// Rate limiter configuration
	RateLimiterConfig rlconfig.RateLimiterConfig
}

var config = &Config{}

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

	err = viper.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		logLevelHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))
	if err != nil {
		fmt.Println(err)
		panic("Error unmarshalling config")
	}

	log.SetLogger(NewLogger(config.LogLevel))

	err = rlconfig.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	config.RateLimiterConfig = *rlconfig.GetConfig()

	log.Log(log.Info, "Config loaded successfully")
}

func logLevelHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Bypass hook if source type is not a string or uint, or the target type is not log.LogLevel
		if (f.Kind() != reflect.String && f.Kind() != reflect.Uint) || t != reflect.TypeOf(log.Debug) {
			return data, nil
		}
		// Check that the data is string
		if f.Kind() == reflect.String {
			// Parse the string into a LogLevel
			logLevel, err := log.ParseLogLevel(data.(string))
			if err != nil {
				return nil, err
			}
			return logLevel, nil
		}
		// Data is uint
		logLevel := log.LogLevel(data.(uint8))
		if logLevel < log.Debug {
			return log.Debug, nil
		} else if logLevel > log.Panic {
			return log.Panic, nil
		}
		return logLevel, nil
	}
}
