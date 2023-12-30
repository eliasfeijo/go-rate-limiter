package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/store"
)

type RateLimiterMiddleware struct {
	store       store.IpStore
	handler     http.Handler
	rateLimiter *limiter.RateLimiter
}

func NewRateLimitMiddleware(config *limiter.RateLimiterConfig, storeStrategy string) *RateLimiterMiddleware {
	if storeStrategy == store.RedisStoreStrategy {
		store.CreateRedisClient()
	}
	return &RateLimiterMiddleware{
		store:       make(store.IpStore),
		handler:     http.DefaultServeMux,
		rateLimiter: limiter.NewRateLimiter(config, storeStrategy),
	}
}

func (m *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	m.handler = next
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.ServeHTTP(w, r)
	})
}

func (m *RateLimiterMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	token := r.Header.Get(m.rateLimiter.Config.TokensHeaderKey)
	if m.rateLimiter.Limit(ip, token) {
		cancelRequest(w)
		return
	}
	m.handler.ServeHTTP(w, r)
}

func cancelRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusTooManyRequests)
	fmt.Fprintf(w, "You have reached the maximum number of requests or actions allowed within a certain time frame")
	return
}
