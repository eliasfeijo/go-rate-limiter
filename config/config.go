package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type RateLimiterTokenConfig struct {
	// Max requests per token
	MaxRequests uint64
	// Token limit duration in seconds (the amount of time the max requests are allowed in)
	LimitInSeconds uint64
	// Token block duration in seconds (the amount of time the token is blocked for after exceeding the max requests)
	BlockInSeconds uint64
}

type RateLimiterRedisConfig struct {
	// Redis host
	Host string `mapstructure:"RATE_LIMITER_REDIS_HOST"`
	// Redis port
	Port string `mapstructure:"RATE_LIMITER_REDIS_PORT"`
	// Redis password
	Password string `mapstructure:"RATE_LIMITER_REDIS_PASSWORD"`
	// Redis database
	DB int `mapstructure:"RATE_LIMITER_REDIS_DB"`
}

// A map of tokens configurations
type TokensConfig map[string]*RateLimiterTokenConfig

type Config struct {
	// Max requests per IP address
	RateLimiterIpAddressMaxRequests uint64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS"`
	// IP Address limit duration in seconds (the amount of time the max requests are allowed in)
	RateLimiterIpAddressLimitInSeconds uint64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS"`
	// IP Address block duration in seconds (the amount of time the IP address is blocked for after exceeding the max requests)
	RateLimiterIpAddressBlockInSeconds uint64 `mapstructure:"RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS"`
	// The requests' Header key to use for the tokens
	RateLimiterTokensHeaderKey string `mapstructure:"RATE_LIMITER_TOKENS_HEADER_KEY"`
	// A list of tokens separated by a comma and their respective max requests, limit and block durations in seconds separated by a colon
	RateLimiterTokensConfigTuple string `mapstructure:"RATE_LIMITER_TOKENS_CONFIG_TUPLE"`
	// The strategy to use for the store
	RateLimiterStoreStrategy string `mapstructure:"RATE_LIMITER_STORE_STRATEGY"`

	// A map of tokens and their respective max requests, limit and block durations in seconds
	TokensConfig

	// Redis configuration
	RateLimiterRedisConfig `mapstructure:",squash"`
}

var config = &Config{}

func LoadConfig() (err error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS", 2)
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS", 1)
	viper.SetDefault("RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS", 5)
	viper.SetDefault("RATE_LIMITER_TOKENS_HEADER_KEY", "API_KEY")
	viper.SetDefault("RATE_LIMITER_TOKENS_CONFIG_TUPLE", "")
	viper.SetDefault("RATE_LIMITER_STORE_STRATEGY", "in_memory")
	viper.SetDefault("RATE_LIMITER_REDIS_HOST", "localhost")
	viper.SetDefault("RATE_LIMITER_REDIS_PORT", "6379")
	viper.SetDefault("RATE_LIMITER_REDIS_PASSWORD", "")
	viper.SetDefault("RATE_LIMITER_REDIS_DB", 0)

	err = viper.ReadInConfig()
	if err != nil {
		if err2, ok := err.(*os.PathError); !ok {
			err = err2
			log.Log(log.Error, "Error reading config file")
			return
		}
	}

	err = viper.Unmarshal(config)

	if config.RateLimiterTokensConfigTuple != "" {
		config.TokensConfig = make(TokensConfig)
		viper.UnmarshalKey("RATE_LIMITER_TOKENS_CONFIG_TUPLE", &config.TokensConfig, viper.DecodeHook(tokensMapHookFunc()))
	}

	log.Log(log.Debug, "IP Address Max Requests:", config.RateLimiterIpAddressMaxRequests)
	log.Log(log.Debug, "IP Address Limit In Seconds:", config.RateLimiterIpAddressLimitInSeconds)
	log.Log(log.Debug, "IP Address Block In Seconds:", config.RateLimiterIpAddressBlockInSeconds)
	for token, tokenConfig := range config.TokensConfig {
		log.Log(log.Debug, "Token:", token)
		log.Log(log.Debug, tokenConfig)
	}
	return
}

func GetConfig() *Config {
	return config
}

func tokensMapHookFunc() mapstructure.DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		// Check that the data is string
		if f.Kind() != reflect.String {
			return data, nil
		}

		// Check that the target type is our custom type
		if t != reflect.TypeOf(TokensConfig{}) {
			return data, nil
		}

		// Parse the tokens config string
		tuples := strings.Split(data.(string), ",")

		// Create a map to store the tokens config
		tokensConfig := make(map[string]*RateLimiterTokenConfig)

		// Loop through the tokens config tuples
		for _, tuple := range tuples {
			// Split the tokens config tuple into token and config
			parsed := strings.Split(tuple, ":")
			// Check that the token config tuple is valid
			if len(parsed) != 4 {
				return nil, fmt.Errorf("Invalid token config tuple: %s", tuple)
			}
			token := parsed[0]
			MaxRequests, err := strconv.ParseUint(parsed[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid token config tuple: %s", tuple)
			}
			LimitInSeconds, err := strconv.ParseUint(parsed[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid token config tuple: %s", tuple)
			}
			BlockInSeconds, err := strconv.ParseUint(parsed[3], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Invalid token config tuple: %s", tuple)
			}
			tokensConfig[token] = &RateLimiterTokenConfig{
				MaxRequests,
				LimitInSeconds,
				BlockInSeconds,
			}
		}

		return tokensConfig, nil
	}
}

func (t *RateLimiterTokenConfig) String() string {
	return fmt.Sprintf(
		"Max Requests: %d, Limit In Seconds: %d, Block In Seconds: %d",
		t.MaxRequests,
		t.LimitInSeconds,
		t.BlockInSeconds,
	)
}
