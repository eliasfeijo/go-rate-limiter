package limiter

import (
	"fmt"
	"sync"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/store"
)

type RateLimiterConfig struct {
	MaxRequestsIpAddress    uint64
	LimitInSecondsIpAddress uint64
	BlockInSecondsIpAddress uint64
	TokensHeaderKey         string
	TokensConfig            config.TokensConfig
}

type RateLimiter struct {
	Config        *RateLimiterConfig
	Store         store.IpStore
	storeStrategy string
	mutex         sync.Mutex
}

func NewRateLimiter(config *RateLimiterConfig, storeStrategy string) *RateLimiter {
	return &RateLimiter{
		Config:        config,
		Store:         make(store.IpStore),
		storeStrategy: storeStrategy,
		mutex:         sync.Mutex{},
	}
}

func (rl *RateLimiter) Limit(ip string, token string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	var s store.Store

	if s, ok := rl.Store[ip][token]; !ok {
		maxRequests := rl.Config.MaxRequestsIpAddress
		limitInSeconds := rl.Config.LimitInSecondsIpAddress
		blockInSeconds := rl.Config.BlockInSecondsIpAddress
		if token != "" {
			maxRequests = rl.Config.TokensConfig[token].MaxRequests
			limitInSeconds = rl.Config.TokensConfig[token].LimitInSeconds
			blockInSeconds = rl.Config.TokensConfig[token].BlockInSeconds
		}
		rl.Store[ip] = make(store.TokenStore)
		switch rl.storeStrategy {
		default:
		case store.InMemoryStoreStrategy:
			rl.Store[ip][token] = store.NewInMemoryStore(uint(maxRequests), limitInSeconds, blockInSeconds)
		}
		s = rl.Store[ip][token]
	} else {
		if s.ShouldRefresh() {
			fmt.Printf("IP: %s, Token: %s, Refreshing", ip, token)
			s.Refresh()
		} else {
			if s.IsBlocked() {
				fmt.Printf("IP: %s, Token: %s, Blocked for more %d seconds\n", ip, token, s.RemainingBlockTime())
				return true
			}
			fmt.Println("Incrementing")
			s.Hit()
		}
	}
	s, _ = rl.Store[ip][token]
	fmt.Printf("\n\nStore:%v\nLastHit:%v\n\n", s, s.LastHit())
	lastHit := uint64(time.Now().Unix() - s.LastHit().Unix())
	fmt.Printf("IP: %s, Token: %s, HitCount: %d, LastHit: %s, Blocked: %v\n", ip, token, s.HitCount(), fmt.Sprint(lastHit)+"s", s.IsBlocked())
	if s.ShouldLimit() {
		fmt.Println("Blocking")
		s.Block()
		return true
	}
	return false
}
