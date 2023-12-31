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
	rl                 *limiter.RateLimiter
	store              store.IpStore
	mockStore          *mocks.MockStore
	onPrepareMockStore OnPrepareMockStore
}

type OnPrepareMockStore func(store *mocks.MockStore) store.Store

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
		suite.mockStore = &mocks.MockStore{}
		if suite.onPrepareMockStore != nil {
			return suite.onPrepareMockStore(suite.mockStore)
		}
		return suite.mockStore
	})
}
func (s *LimiterTestSuite) TestShouldCreateStoreWhenItDoesNotExist() {
	s.onPrepareMockStore = func(store *mocks.MockStore) store.Store {
		store.On("HitCount").Return(uint(1))
		store.On("IsBlocked").Return(false)
		store.On("RemainingBlockTime").Return(uint(0))
		store.On("ShouldLimit").Return(false)
		store.On("ShouldRefresh").Return(false)
		store.On("LastHit").Return(uint(0))
		return store
	}
	ip := "123"
	s.rl.Limit(ip, "")
	store := s.rl.Store[ip][""]
	assert.NotNil(s.T(), store)
	assert.Equal(s.T(), uint(1), store.HitCount())
	assert.Equal(s.T(), false, store.IsBlocked())
	assert.Equal(s.T(), uint(0), store.RemainingBlockTime())
	assert.Equal(s.T(), int64(0), store.LastHit().Unix())
	s.mockStore.AssertCalled(s.T(), "ShouldLimit")
}

func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
