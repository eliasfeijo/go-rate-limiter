package limiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/eliasfeijo/go-rate-limiter/store"
)

type RateLimiter struct {
	Config *config.RateLimiterConfig
	Store  store.IpStore
	mutex  sync.Mutex
}

func NewRateLimiter(config *config.RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		Config: config,
		Store:  make(store.IpStore),
		mutex:  sync.Mutex{},
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
		switch rl.Config.StoreStrategy {
		case store.RedisStoreStrategy:
			rl.Store[ip][token] = store.NewRedisStore(ip, token, uint(maxRequests), limitInSeconds, blockInSeconds)
		default:
		case store.InMemoryStoreStrategy:
			rl.Store[ip][token] = store.NewInMemoryStore(uint(maxRequests), limitInSeconds, blockInSeconds)
		}
		s = rl.Store[ip][token]
	} else {
		if s.ShouldRefresh() {
			log.Logf(log.Debug, "IP: %s, Token: %s, Refreshing", ip, token)
			s.Refresh()
		} else {
			if s.IsBlocked() {
				log.Logf(log.Debug, "IP: %s, Token: %s, Blocked for more %d seconds", ip, token, s.RemainingBlockTime())
				return true
			}
			log.Log(log.Debug, "Incrementing")
			s.Hit()
		}
	}
	s, _ = rl.Store[ip][token]
	lastHit := time.Now().Unix() - s.LastHit().Unix()
	log.Logf(log.Debug, "IP: %s, Token: %s, HitCount: %d, LastHit: %s, Blocked: %v\n", ip, token, s.HitCount(), fmt.Sprint(lastHit)+"s", s.IsBlocked())
	if s.ShouldLimit() {
		log.Log(log.Debug, "Blocking")
		s.Block()
		return true
	}
	return false
}
