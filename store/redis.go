package store

import (
	"context"
	"strconv"
	"time"

	"github.com/eliasfeijo/go-rate-limiter/config"
	"github.com/eliasfeijo/go-rate-limiter/log"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type RedisStore struct {
	config *StoreConfig
	key    string
	ctx    context.Context
}

func CreateRedisClient() {
	if rdb != nil {
		return
	}
	cfg := config.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Host + ":" + cfg.RedisConfig.Port,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	log.Log(log.Info, "Redis connection created successfully: ", rdb)
}

func NewRedisStore(ip string, token string, config *StoreConfig) *RedisStore {
	ctx := context.Background()
	key := ip + ":" + token
	rdb.Set(ctx, key+":hitCount", 1, 0)
	rdb.Set(ctx, key+":lastHit", time.Now().Unix(), 0)
	rdb.Set(ctx, key+":isBlocked", false, 0)
	return &RedisStore{config, key, ctx}
}

func (s *RedisStore) ShouldLimit() bool {
	hitCountString, err := rdb.Get(s.ctx, s.key+":hitCount").Result()
	if err != nil {
		hitCountString = "0"
	}
	hitCount, err := strconv.ParseUint(hitCountString, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(hitCount) > s.config.MaxRequests
}

func (s *RedisStore) ShouldRefresh() bool {
	isBlocked, _ := rdb.Get(s.ctx, s.key+":isBlocked").Result()
	if isBlocked == "true" {
		return true
	}

	lastHitString, err := rdb.Get(s.ctx, s.key+":lastHit").Result()
	if err != nil {
		lastHitString = "0"
	}
	lastHit, err := strconv.ParseInt(lastHitString, 10, 64)
	if err != nil {
		panic(err)
	}

	return uint(time.Now().Unix()-lastHit) > s.config.LimitInSeconds
}

func (s *RedisStore) Refresh() {
	rdb.Set(s.ctx, s.key+":hitCount", 1, 0)
	rdb.Set(s.ctx, s.key+":lastHit", time.Now().Unix(), 0)
	rdb.Set(s.ctx, s.key+":isBlocked", false, 0)
}

func (s *RedisStore) IsBlocked() bool {
	isBlockedString, err := rdb.Get(s.ctx, s.key+":isBlocked").Result()
	if err != nil {
		isBlockedString = "false"
	}
	isBlocked, err := strconv.ParseBool(isBlockedString)
	if err != nil {
		panic(err)
	}
	return isBlocked
}

func (s *RedisStore) RemainingBlockTime() uint {
	expDuration, err := rdb.ExpireTime(s.ctx, s.key+":isBlocked").Result()
	if err != nil {
		expDuration = time.Duration(0) * time.Second
	}
	exp := expDuration.Seconds()
	return uint(exp) - uint(time.Now().Unix())
}

func (s *RedisStore) Block() {
	rdb.Set(s.ctx, s.key+":isBlocked", true, time.Duration(s.config.BlockInSeconds)*time.Second)
}

func (s *RedisStore) Hit() {
	rdb.Incr(s.ctx, s.key+":hitCount")
	rdb.Set(s.ctx, s.key+":lastHit", time.Now().Unix(), 0)
	if s.ShouldLimit() {
		s.Block()
	}
}

func (s *RedisStore) LastHit() time.Time {
	lastHitString, err := rdb.Get(s.ctx, s.key+":lastHit").Result()
	if err != nil {
		lastHitString = "0"
	}
	lastHit, err := strconv.ParseInt(lastHitString, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(lastHit, 0)
}

func (s *RedisStore) HitCount() uint {
	hitCountString, err := rdb.Get(s.ctx, s.key+":hitCount").Result()
	if err != nil {
		hitCountString = "0"
	}
	hitCount, err := strconv.ParseUint(hitCountString, 10, 64)
	if err != nil {
		panic(err)
	}
	return uint(hitCount)
}
