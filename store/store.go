package store

import "time"

const (
	InMemoryStoreStrategy = "in_memory"
	RedisStoreStrategy    = "redis"
)

type Store interface {
	ShouldLimit() bool
	ShouldRefresh() bool
	Refresh()
	IsBlocked() bool
	RemainingBlockTime() uint64
	Block()
	Hit()
	LastHit() time.Time
	HitCount() uint
}

type TokenStore map[string]Store
type IpStore map[string]TokenStore
