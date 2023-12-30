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
	LastHit        time.Time
	Blocked        bool
	LimitInSeconds int64
	BlockInSeconds int64
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
			store[ip] = &Store{HitCount: 1, LastHit: time.Now(), Blocked: false, LimitInSeconds: limitInSeconds, BlockInSeconds: blockInSeconds}
		} else {
			if store[ip].shouldRefresh() {
				fmt.Println("Refreshing")
				store[ip].HitCount = 1
				store[ip].LastHit = time.Now()
				store[ip].Blocked = false
			} else {
				if store[ip].Blocked {
					fmt.Println("Blocked")
					limit(w)
					return
				}
				fmt.Println("Incrementing")
				store[ip].HitCount++
				store[ip].LastHit = time.Now()
			}
		}
		fmt.Printf("IP: %s, HitCount: %d, LastHit: %s, BlockedAt: %v\n", ip, store[ip].HitCount, store[ip].LastHit, store[ip].Blocked)
		if store[ip].shouldLimit() {
			fmt.Println("Limiting")
			store[ip].Blocked = true
			limit(w)
			return
		}
		fmt.Fprintf(w, "Hello, %s! You have hit me %d times.", ip, store[ip].HitCount)
		w.Write([]byte("Welcome!"))
	})
	http.ListenAndServe(":"+cfg.Port, r)
}

func (s *Store) shouldRefresh() bool {
	LastHit := time.Now().Unix() - s.LastHit.Unix()
	if s.Blocked {
		return LastHit > s.BlockInSeconds
	}
	return LastHit > s.LimitInSeconds
}

func (s *Store) shouldLimit() bool {
	return s.HitCount > 2
}

func limit(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
	return
}
