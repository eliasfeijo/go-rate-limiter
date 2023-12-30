package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var mutex = &sync.Mutex{}

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	cfg := config.GetConfig()
	maxRequestsIpAddress := cfg.RateLimiterIpAddressMaxRequests
	limitInSecondsIpAddress := cfg.RateLimiterIpAddressLimitInSeconds
	blockInSecondsIpAddress := cfg.RateLimiterIpAddressBlockInSeconds
	tokensHeaderKey := cfg.RateLimiterTokensHeaderKey

	store := make(limiter.IpStore)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		ip := strings.Split(r.RemoteAddr, ":")[0]
		token := r.Header.Get(tokensHeaderKey)
		if _, ok := store[ip][token]; !ok {
			maxRequests := maxRequestsIpAddress
			limitInSeconds := limitInSecondsIpAddress
			blockInSeconds := blockInSecondsIpAddress
			if token != "" {
				maxRequests = cfg.TokensConfig[token].MaxRequests
				limitInSeconds = cfg.TokensConfig[token].LimitInSeconds
				blockInSeconds = cfg.TokensConfig[token].BlockInSeconds
			}
			store[ip] = make(map[string]*limiter.Store)
			store[ip][token] = &limiter.Store{
				HitCount:       1,
				LastHit:        time.Now(),
				Blocked:        false,
				MaxRequests:    uint(maxRequests),
				LimitInSeconds: limitInSeconds,
				BlockInSeconds: blockInSeconds,
			}
		} else {
			if store[ip][token].ShouldRefresh() {
				fmt.Printf("IP: %s, Token: %s, Refreshing", ip, token)
				store[ip][token].HitCount = 1
				store[ip][token].LastHit = time.Now()
				store[ip][token].Blocked = false
			} else {
				if store[ip][token].Blocked {
					remainingTime := store[ip][token].BlockInSeconds - uint64(time.Now().Unix()-store[ip][token].LastHit.Unix())
					fmt.Printf("IP: %s, Token: %s, Blocked for more %d seconds\n", ip, token, remainingTime)
					limit(w)
					return
				}
				fmt.Println("Incrementing")
				store[ip][token].HitCount++
				store[ip][token].LastHit = time.Now()
			}
		}
		lastHit := uint(time.Now().Unix() - store[ip][token].LastHit.Unix())
		fmt.Printf("IP: %s, Token: %s, HitCount: %d, LastHit: %s, Blocked: %v\n", ip, token, store[ip][token].HitCount, fmt.Sprint(lastHit)+"s", store[ip][token].Blocked)
		if store[ip][token].ShouldLimit() {
			fmt.Println("Limiting")
			store[ip][token].Blocked = true
			limit(w)
			return
		}
		fmt.Fprintf(w, "Hello, %s! You have hit me %d times.", ip, store[ip][token].HitCount)
	})
	http.ListenAndServe(":"+cfg.Port, r)
}

func limit(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
	return
}
