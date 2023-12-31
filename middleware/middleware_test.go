package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/go-chi/chi/v5"
)

func testRequest(t *testing.T, server *httptest.Server) (*http.Response, string) {
	// Create a mock HTTP request
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRateLimiterMiddleware_ServeHTTP_Success(t *testing.T) {
	cfg := &config.RateLimiterConfig{
		IpAddressMaxRequests:    3,
		IpAddressLimitInSeconds: 1,
		IpAddressBlockInSeconds: 5,
		MapTokenConfig:          nil,
		TokensHeaderKey:         "API_KEY",
		StoreStrategy:           "in_memory",
		RedisConfig:             config.RedisConfig{},
	}

	// Create a new instance of the RateLimiterMiddleware
	middleware := NewRateLimitMiddleware(cfg)

	r := chi.NewRouter()
	r.Use(middleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request accepted"))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	// Create a mock HTTP request
	if resp, body := testRequest(t, server); body != "root" && resp.StatusCode != 200 {
		t.Fatalf(body)
	}
}

func TestRateLimiterMiddleware_ServeHTTP_TooManyRequests(t *testing.T) {

	cfg := &config.RateLimiterConfig{
		IpAddressMaxRequests:    1,
		IpAddressLimitInSeconds: 1,
		IpAddressBlockInSeconds: 1,
		MapTokenConfig:          nil,
		TokensHeaderKey:         "API_KEY",
		StoreStrategy:           "in_memory",
		RedisConfig:             config.RedisConfig{},
	}

	// Create a new instance of the RateLimiterMiddleware
	middleware := NewRateLimitMiddleware(cfg)

	r := chi.NewRouter()
	r.Use(middleware.Handler)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request accepted"))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	// Create a mock HTTP request
	testRequest(t, server)
	testRequest(t, server)
	if resp, _ := testRequest(t, server); resp.StatusCode != 429 {
		t.Fatalf("Expected status code 429, got %d", resp.StatusCode)
	}
}
