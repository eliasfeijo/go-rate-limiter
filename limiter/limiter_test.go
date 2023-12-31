package limiter_test

// Basic imports
import (
	"testing"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/limiter"
	"github.com/eliasfeijo/go-rate-limiter/mocks"
	"github.com/eliasfeijo/go-rate-limiter/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LimiterTestSuite struct {
	suite.Suite
	rateLimiterConfig *config.RateLimiterConfig
}

var ip = "123"

var defaultRateLimiterConfig = &config.RateLimiterConfig{
	IpAddressMaxRequests:    3,
	IpAddressLimitInSeconds: 1,
	IpAddressBlockInSeconds: 5,
	MapTokenConfig:          nil,
	TokensHeaderKey:         "API_KEY",
	StoreStrategy:           "test",
	RedisConfig:             config.RedisConfig{},
}

func (suite *LimiterTestSuite) SetupTest() {
	suite.rateLimiterConfig = defaultRateLimiterConfig
}

func (s *LimiterTestSuite) TestShouldCreateStoreWhenItDoesNotExist() {
	rl := limiter.NewRateLimiter(s.rateLimiterConfig, make(store.IpStore), func(store store.Store) store.Store {
		mockStore := &mocks.MockStore{}
		mockStore.On("ShouldLimit").Return(false)
		return mockStore
	})
	result := rl.Limit(ip, "")
	assert.False(s.T(), result)
	assert.NotNil(s.T(), rl.Store[ip])
	assert.NotNil(s.T(), rl.Store[ip][""])
	assert.IsType(s.T(), &mocks.MockStore{}, rl.Store[ip][""])
}

func (s *LimiterTestSuite) TestRefresh() {
	mockStore := &mocks.MockStore{}
	mockStore.On("ShouldRefresh").Return(true)
	mockStore.On("Refresh").Return()
	mockStore.On("IsBlocked").Return(false)
	mockStore.On("Hit").Return()
	mockStore.On("ShouldLimit").Return(false)
	mockStore.On("Block").Return()
	ipStore := make(store.IpStore)
	ipStore[ip] = make(store.TokenStore)
	ipStore[ip][""] = mockStore
	rl := limiter.NewRateLimiter(s.rateLimiterConfig, ipStore, nil)
	result := rl.Limit(ip, "")
	assert.False(s.T(), result)
	store := rl.Store[ip][""].(*mocks.MockStore)
	store.AssertCalled(s.T(), "ShouldRefresh")
	store.AssertCalled(s.T(), "Refresh")
	store.AssertNotCalled(s.T(), "IsBlocked")
	store.AssertNotCalled(s.T(), "Hit")
	store.AssertCalled(s.T(), "ShouldLimit")
	store.AssertNotCalled(s.T(), "Block")
}

func (s *LimiterTestSuite) TestIsBlocked() {
	mockStore := &mocks.MockStore{}
	mockStore.On("ShouldRefresh").Return(false)
	mockStore.On("IsBlocked").Return(true)
	mockStore.On("Hit").Return()
	mockStore.On("ShouldLimit").Return(false)
	mockStore.On("Block").Return()
	ipStore := make(store.IpStore)
	ipStore[ip] = make(store.TokenStore)
	ipStore[ip][""] = mockStore
	rl := limiter.NewRateLimiter(s.rateLimiterConfig, ipStore, nil)
	result := rl.Limit(ip, "")
	assert.True(s.T(), result)
	store := rl.Store[ip][""].(*mocks.MockStore)
	store.AssertCalled(s.T(), "ShouldRefresh")
	store.AssertCalled(s.T(), "IsBlocked")
	store.AssertNotCalled(s.T(), "Hit")
	store.AssertNotCalled(s.T(), "ShouldLimit")
	store.AssertNotCalled(s.T(), "Block")
}

func (s *LimiterTestSuite) TestShouldLimitAndBlock() {
	mockStore := &mocks.MockStore{}
	mockStore.On("ShouldRefresh").Return(false)
	mockStore.On("IsBlocked").Return(false)
	mockStore.On("Hit").Return()
	mockStore.On("ShouldLimit").Return(true)
	mockStore.On("Block").Return()
	ipStore := make(store.IpStore)
	ipStore[ip] = make(store.TokenStore)
	ipStore[ip][""] = mockStore
	rl := limiter.NewRateLimiter(s.rateLimiterConfig, ipStore, nil)
	result := rl.Limit(ip, "")
	assert.True(s.T(), result)
	store := rl.Store[ip][""].(*mocks.MockStore)
	store.AssertCalled(s.T(), "ShouldRefresh")
	store.AssertCalled(s.T(), "IsBlocked")
	store.AssertCalled(s.T(), "Hit")
	store.AssertCalled(s.T(), "ShouldLimit")
	store.AssertCalled(s.T(), "Block")
}

func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
