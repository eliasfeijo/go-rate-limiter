package limiter

import (
	"fmt"
	"sync"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/mocks"
	"github.com/eliasfeijo/go-rate-limiter/store"
)

type RateLimiter struct {
	Config         *config.RateLimiterConfig
	Store          store.IpStore
	mutex          sync.Mutex
	onStoreCreated store.StoreCreatedCallback
}

func NewRateLimiter(config *config.RateLimiterConfig, storeCreatedCallback store.StoreCreatedCallback) *RateLimiter {
	return &RateLimiter{
		Config:         config,
		Store:          make(store.IpStore),
		mutex:          sync.Mutex{},
		onStoreCreated: storeCreatedCallback,
	}
}

func (rl *RateLimiter) Limit(ip string, token string) bool {
	// rl.mutex.Lock()
	// defer rl.mutex.Unlock()

	var s store.Store

	if s, ok := rl.Store[ip][token]; !ok {
		log.Log(log.Debug, "Creating new store")
		maxRequests := rl.Config.IpAddressMaxRequests
		limitInSeconds := rl.Config.IpAddressLimitInSeconds
		blockInSeconds := rl.Config.IpAddressBlockInSeconds
		if token != "" {
			maxRequests = rl.Config.MapTokenConfig[token].MaxRequests
			limitInSeconds = rl.Config.MapTokenConfig[token].LimitInSeconds
			blockInSeconds = rl.Config.MapTokenConfig[token].BlockInSeconds
		}
		rl.Store[ip] = make(store.TokenStore)
		storeConfig := &store.StoreConfig{
			MaxRequests:    maxRequests,
			LimitInSeconds: limitInSeconds,
			BlockInSeconds: blockInSeconds,
		}
		switch rl.Config.StoreStrategy {
		case "test":
		case "mock":
			s = mocks.NewMockStore()
			fmt.Println(s)
		case store.RedisStoreStrategy:
			s = store.NewRedisStore(ip, token, storeConfig)
		default:
		case store.InMemoryStoreStrategy:
			s = store.NewInMemoryStore(storeConfig)
		}
		if rl.onStoreCreated != nil {
			s = rl.onStoreCreated(s)
		}
		rl.Store[ip][token] = s
	} else {
		if s.ShouldRefresh() {
			log.Logf(log.Debug, "IP: %s, Token: %s, Refreshing", ip, token)
			s.Refresh()
		} else {
			if s.IsBlocked() {
				return true
			}
			log.Log(log.Debug, "Incrementing")
			s.Hit()
		}
	}
	s, _ = rl.Store[ip][token]
	if s.ShouldLimit() {
		log.Log(log.Debug, "Blocking")
		s.Block()
		return true
	}
	return false
}
