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

type TokenConfig struct {
	// Max requests per token
	MaxRequests uint
	// Token limit duration in seconds (the amount of time the max requests are allowed in)
	LimitInSeconds uint
	// Token block duration in seconds (the amount of time the token is blocked for after exceeding the max requests)
	BlockInSeconds uint
}

// A map of tokens configurations
type MapTokenConfig map[string]*TokenConfig

type RedisConfig struct {
	// Redis host
	Host string `mapstructure:"RATE_LIMITER_REDIS_HOST"`
	// Redis port
	Port string `mapstructure:"RATE_LIMITER_REDIS_PORT"`
	// Redis password
	Password string `mapstructure:"RATE_LIMITER_REDIS_PASSWORD"`
	// Redis database
	DB int `mapstructure:"RATE_LIMITER_REDIS_DB"`
}

type RateLimiterConfig struct {
	// Max requests per IP address
	IpAddressMaxRequests uint `mapstructure:"RATE_LIMITER_IP_ADDRESS_MAX_REQUESTS"`
	// IP Address limit duration in seconds (the amount of time the max requests are allowed in)
	IpAddressLimitInSeconds uint `mapstructure:"RATE_LIMITER_IP_ADDRESS_LIMIT_IN_SECONDS"`
	// IP Address block duration in seconds (the amount of time the IP address is blocked for after exceeding the max requests)
	IpAddressBlockInSeconds uint `mapstructure:"RATE_LIMITER_IP_ADDRESS_BLOCK_IN_SECONDS"`
	// The requests' Header key to use for the tokens
	TokensHeaderKey string `mapstructure:"RATE_LIMITER_TOKENS_HEADER_KEY"`
	// A list of tokens separated by a comma and their respective max requests, limit and block durations in seconds separated by a colon
	MapTokenConfigTuple string `mapstructure:"RATE_LIMITER_TOKENS_CONFIG_TUPLE"`
	// The strategy to use for the store
	StoreStrategy string `mapstructure:"RATE_LIMITER_STORE_STRATEGY"`

	// A map of tokens and their respective max requests, limit and block durations in seconds
	MapTokenConfig `mapstructure:"RATE_LIMITER_TOKENS_CONFIG_TUPLE"`

	// Redis configuration
	RedisConfig `mapstructure:",squash"`
}

var config = &RateLimiterConfig{}

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

	err = viper.Unmarshal(config, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		tokensMapHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	)))

	log.Log(log.Debug, "IP Address Max Requests:", config.IpAddressMaxRequests)
	log.Log(log.Debug, "IP Address Limit In Seconds:", config.IpAddressLimitInSeconds)
	log.Log(log.Debug, "IP Address Block In Seconds:", config.IpAddressBlockInSeconds)
	for token, tokenConfig := range config.MapTokenConfig {
		log.Log(log.Debug, "Token:", token)
		log.Log(log.Debug, tokenConfig)
	}
	return
}

func GetConfig() *RateLimiterConfig {
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
		if t != reflect.TypeOf(MapTokenConfig{}) {
			return data, nil
		}

		// Parse the tokens config string
		tuples := strings.Split(data.(string), ",")

		// Create a map to store the tokens config
		MapTokenConfig := make(map[string]*TokenConfig)

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
			MapTokenConfig[token] = &TokenConfig{
				MaxRequests:    uint(MaxRequests),
				LimitInSeconds: uint(LimitInSeconds),
				BlockInSeconds: uint(BlockInSeconds),
			}
		}

		return MapTokenConfig, nil
	}
}

func (t *TokenConfig) String() string {
	return fmt.Sprintf(
		"Max Requests: %d, Limit In Seconds: %d, Block In Seconds: %d",
		t.MaxRequests,
		t.LimitInSeconds,
		t.BlockInSeconds,
	)
}
