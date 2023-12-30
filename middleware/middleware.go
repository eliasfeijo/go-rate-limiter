package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
)

type RateLimiterConfig struct {
	MaxRequestsIpAddress    uint64
	LimitInSecondsIpAddress uint64
	BlockInSecondsIpAddress uint64
	TokensHeaderKey         string
	TokensConfig            config.TokensConfig
}

type RateLimiter struct {
	store   limiter.IpStore
	config  *RateLimiterConfig
	mutex   sync.Mutex
	handler http.Handler
}

func NewRateLimitMiddleware(config *RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		store:   make(limiter.IpStore),
		config:  config,
		mutex:   sync.Mutex{},
		handler: http.DefaultServeMux,
	}
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	rl.handler = next
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	ip := strings.Split(r.RemoteAddr, ":")[0]
	token := r.Header.Get(rl.config.TokensHeaderKey)
	if _, ok := rl.store[ip][token]; !ok {
		maxRequests := rl.config.MaxRequestsIpAddress
		limitInSeconds := rl.config.LimitInSecondsIpAddress
		blockInSeconds := rl.config.BlockInSecondsIpAddress
		if token != "" {
			maxRequests = rl.config.TokensConfig[token].MaxRequests
			limitInSeconds = rl.config.TokensConfig[token].LimitInSeconds
			blockInSeconds = rl.config.TokensConfig[token].BlockInSeconds
		}
		rl.store[ip] = make(map[string]*limiter.Store)
		rl.store[ip][token] = &limiter.Store{
			HitCount:       1,
			LastHit:        time.Now(),
			Blocked:        false,
			MaxRequests:    uint(maxRequests),
			LimitInSeconds: limitInSeconds,
			BlockInSeconds: blockInSeconds,
		}
	} else {
		if rl.store[ip][token].ShouldRefresh() {
			fmt.Printf("IP: %s, Token: %s, Refreshing", ip, token)
			rl.store[ip][token].HitCount = 1
			rl.store[ip][token].LastHit = time.Now()
			rl.store[ip][token].Blocked = false
		} else {
			if rl.store[ip][token].Blocked {
				remainingTime := rl.store[ip][token].BlockInSeconds - uint64(time.Now().Unix()-rl.store[ip][token].LastHit.Unix())
				fmt.Printf("IP: %s, Token: %s, Blocked for more %d seconds\n", ip, token, remainingTime)
				limit(w)
				return
			}
			fmt.Println("Incrementing")
			rl.store[ip][token].HitCount++
			rl.store[ip][token].LastHit = time.Now()
		}
	}
	lastHit := uint(time.Now().Unix() - rl.store[ip][token].LastHit.Unix())
	fmt.Printf("IP: %s, Token: %s, HitCount: %d, LastHit: %s, Blocked: %v\n", ip, token, rl.store[ip][token].HitCount, fmt.Sprint(lastHit)+"s", rl.store[ip][token].Blocked)
	if rl.store[ip][token].ShouldLimit() {
		fmt.Println("Limiting")
		rl.store[ip][token].Blocked = true
		limit(w)
		return
	}
	rl.handler.ServeHTTP(w, r)
}

func limit(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
	return
}
