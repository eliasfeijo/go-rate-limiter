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
	Config *RateLimiterConfig
	Store  store.IpStore
	mutex  sync.Mutex
}

func NewRateLimiter(config *RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		Config: config,
		Store:  make(store.IpStore),
		mutex:  sync.Mutex{},
	}
}

func (rl *RateLimiter) Limit(ip string, token string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if _, ok := rl.Store[ip][token]; !ok {
		maxRequests := rl.Config.MaxRequestsIpAddress
		limitInSeconds := rl.Config.LimitInSecondsIpAddress
		blockInSeconds := rl.Config.BlockInSecondsIpAddress
		if token != "" {
			maxRequests = rl.Config.TokensConfig[token].MaxRequests
			limitInSeconds = rl.Config.TokensConfig[token].LimitInSeconds
			blockInSeconds = rl.Config.TokensConfig[token].BlockInSeconds
		}
		rl.Store[ip] = make(map[string]*store.Store)
		rl.Store[ip][token] = &store.Store{
			HitCount:       1,
			LastHit:        time.Now(),
			Blocked:        false,
			MaxRequests:    uint(maxRequests),
			LimitInSeconds: limitInSeconds,
			BlockInSeconds: blockInSeconds,
		}
	} else {
		if rl.Store[ip][token].ShouldRefresh() {
			fmt.Printf("IP: %s, Token: %s, Refreshing", ip, token)
			rl.Store[ip][token].HitCount = 1
			rl.Store[ip][token].LastHit = time.Now()
			rl.Store[ip][token].Blocked = false
		} else {
			if rl.Store[ip][token].Blocked {
				remainingTime := rl.Store[ip][token].BlockInSeconds - uint64(time.Now().Unix()-rl.Store[ip][token].LastHit.Unix())
				fmt.Printf("IP: %s, Token: %s, Blocked for more %d seconds\n", ip, token, remainingTime)
				return true
			}
			fmt.Println("Incrementing")
			rl.Store[ip][token].HitCount++
			rl.Store[ip][token].LastHit = time.Now()
		}
	}
	lastHit := uint(time.Now().Unix() - rl.Store[ip][token].LastHit.Unix())
	fmt.Printf("IP: %s, Token: %s, HitCount: %d, LastHit: %s, Blocked: %v\n", ip, token, rl.Store[ip][token].HitCount, fmt.Sprint(lastHit)+"s", rl.Store[ip][token].Blocked)
	if rl.Store[ip][token].ShouldLimit() {
		fmt.Println("Limiting")
		rl.Store[ip][token].Blocked = true
		return true
	}
	return false
}
