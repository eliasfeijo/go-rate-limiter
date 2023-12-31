package store

import (
	"time"
)

const (
	InMemoryStoreStrategy = "in_memory"
	RedisStoreStrategy    = "redis"
)

type Store interface {
	ShouldLimit() bool
	ShouldRefresh() bool
	Refresh()
	IsBlocked() bool
	RemainingBlockTime() uint
	Block()
	Hit()
	LastHit() time.Time
	HitCount() uint
}

type TokenStore map[string]Store
type IpStore map[string]TokenStore

type StoreConfig struct {
	MaxRequests    uint
	LimitInSeconds uint
	BlockInSeconds uint
}

type StoreCreatedCallback func(store Store) Store

// func NewStore(storeStrategy string, ip string, token string, config *StoreConfig) Store {
// 	switch storeStrategy {
// 	case "test":
// 	case "mock":
// 		return mocks.NewMockStore()
// 	case RedisStoreStrategy:
// 		return NewRedisStore(ip, token, config)
// 	default:
// 	case InMemoryStoreStrategy:
// 		return NewInMemoryStore(config)
// 	}
// 	return nil
// }
