package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Store struct {
	HitCount       int
	lastHit        time.Time
	blocked        bool
	limitInSeconds int64
	blockInSeconds int64
}

var store map[string]*Store

func main() {
	err := config.LoadConfig()
	if err != nil {
		panic("Error loading config")
	}

	cfg := config.GetConfig()
	limitInSeconds := cfg.RateLimiterIpAddressLimitInSeconds
	blockInSeconds := cfg.RateLimiterIpAddressBlockInSeconds

	store = make(map[string]*Store)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ip := strings.Split(r.RemoteAddr, ":")[0]
		if _, ok := store[ip]; !ok {
			store[ip] = &Store{HitCount: 1, lastHit: time.Now(), blocked: false, limitInSeconds: limitInSeconds, blockInSeconds: blockInSeconds}
		} else {
			if store[ip].shouldRefresh() {
				fmt.Println("Refreshing")
				store[ip].HitCount = 1
				store[ip].lastHit = time.Now()
				store[ip].blocked = false
			} else {
				if store[ip].blocked {
					fmt.Println("Blocked")
					limit(w)
					return
				}
				fmt.Println("Incrementing")
				store[ip].HitCount++
				store[ip].lastHit = time.Now()
			}
		}
		fmt.Printf("IP: %s, HitCount: %d, LastHit: %s, blockedAt: %v\n", ip, store[ip].HitCount, store[ip].lastHit, store[ip].blocked)
		if store[ip].shouldLimit() {
			fmt.Println("Limiting")
			store[ip].blocked = true
			limit(w)
			return
		}
		fmt.Fprintf(w, "Hello, %s! You have hit me %d times.", ip, store[ip].HitCount)
		w.Write([]byte("Welcome!"))
	})
	http.ListenAndServe(":"+cfg.Port, r)
}

func (s *Store) shouldRefresh() bool {
	lastHit := time.Now().Unix() - s.lastHit.Unix()
	if s.blocked {
		return lastHit > s.blockInSeconds
	}
	return lastHit > s.limitInSeconds
}

func (s *Store) shouldLimit() bool {
	return s.HitCount > 2
}

func limit(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
	return
}
