package store

import "time"

type InMemoryStore struct {
	config    *StoreConfig
	hitCount  uint
	lastHit   time.Time
	isBlocked bool
}

func NewInMemoryStore(config *StoreConfig) *InMemoryStore {
	return &InMemoryStore{
		config:    config,
		hitCount:  1,
		lastHit:   time.Now(),
		isBlocked: false,
	}
}

func (s *InMemoryStore) ShouldLimit() bool {
	return s.hitCount > s.config.MaxRequests
}

func (s *InMemoryStore) ShouldRefresh() bool {
	LastHit := uint64(time.Now().Unix() - s.lastHit.Unix())
	if s.isBlocked {
		return LastHit > s.config.BlockInSeconds
	}
	return LastHit > s.config.LimitInSeconds
}

func (s *InMemoryStore) Refresh() {
	s.hitCount = 1
	s.lastHit = time.Now()
	s.isBlocked = false
}

func (s *InMemoryStore) IsBlocked() bool {
	return s.isBlocked
}

func (s *InMemoryStore) RemainingBlockTime() uint64 {
	return s.config.BlockInSeconds - uint64(time.Now().Unix()-s.lastHit.Unix())
}

func (s *InMemoryStore) Block() {
	s.isBlocked = true
}

func (s *InMemoryStore) Hit() {
	s.hitCount++
	s.lastHit = time.Now()
	if s.ShouldLimit() {
		s.isBlocked = true
	}
}

func (s *InMemoryStore) LastHit() time.Time {
	return s.lastHit
}

func (s *InMemoryStore) HitCount() uint {
	return s.hitCount
}
