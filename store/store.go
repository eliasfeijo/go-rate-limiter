package store

import "time"

type Store struct {
	HitCount       uint
	LastHit        time.Time
	Blocked        bool
	MaxRequests    uint
	LimitInSeconds uint64
	BlockInSeconds uint64
}

type TokenStore map[string]*Store
type IpStore map[string]TokenStore

func (s *Store) ShouldRefresh() bool {
	LastHit := uint64(time.Now().Unix() - s.LastHit.Unix())
	if s.Blocked {
		return LastHit > s.BlockInSeconds
	}
	return LastHit > s.LimitInSeconds
}

func (s *Store) ShouldLimit() bool {
	return s.HitCount > s.MaxRequests
}
