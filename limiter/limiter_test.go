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
	rl        *limiter.RateLimiter
	store     store.IpStore
	mockStore mocks.MockStore
}

func (suite *LimiterTestSuite) SetupTest() {
	suite.store = make(store.IpStore)
	mapTokenConfig := make(config.MapTokenConfig)
	mapTokenConfig["test"] = &config.TokenConfig{
		MaxRequests:    10,
		LimitInSeconds: 60,
		BlockInSeconds: 60,
	}
	suite.rl = limiter.NewRateLimiter(&config.RateLimiterConfig{
		IpAddressMaxRequests:    10,
		IpAddressLimitInSeconds: 60,
		IpAddressBlockInSeconds: 60,
		MapTokenConfig:          mapTokenConfig,
		TokensHeaderKey:         "API_KEY",
		StoreStrategy:           "test",
		RedisConfig: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}, func(store store.Store) store.Store {
		suite.mockStore = mocks.MockStore{}
		suite.mockStore.On("ShouldLimit").Return(false)
		suite.mockStore.On("ShouldRefresh").Return(false)
		suite.mockStore.On("HitCount").Return(uint(1))
		suite.mockStore.On("IsBlocked").Return(false)
		suite.mockStore.On("RemainingBlockTime").Return(uint64(0))
		suite.mockStore.On("LastHit").Return(uint64(0))
		return &suite.mockStore
	})
}
func (s *LimiterTestSuite) TestShouldCreateStoreWhenTokenDoesNotExist() {
	ip := "123"
	s.rl.Limit(ip, "")
	assert.NotNil(s.T(), s.rl.Store[ip][""])
	assert.Equal(s.T(), uint(1), s.rl.Store[ip][""].HitCount())
	assert.Equal(s.T(), false, s.rl.Store[ip][""].IsBlocked())
	assert.Equal(s.T(), uint64(0), s.rl.Store[ip][""].RemainingBlockTime())
	assert.Equal(s.T(), int64(0), s.rl.Store[ip][""].LastHit().Unix())
}

func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
