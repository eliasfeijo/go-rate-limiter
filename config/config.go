package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	// Server port
	Port string `mapstructure:"PORT"`
	// Max requests per IP address
	RateLimiterIpAddressMaxRequests int64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS"`
	// Limit duration in seconds (the amount of time the max requests are allowed in)
	RateLimiterIpAddressLimitInSeconds int64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS"`
	// Block duration in seconds (the amount of time the IP address is blocked for after exceeding the max requests)
	RateLimiterIpAddressBlockInSeconds int64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS"`
}

var config = &Config{}

func LoadConfig() (err error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS", 2)
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS", 1)
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS", 5)

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			fmt.Println("Error reading config file")
			return
		}
	}

	err = viper.Unmarshal(config)

	return
}

func GetConfig() *Config {
	return config
}
