package store

import (
	"testing"
	"time"
)

func TestInMemoryStore_ShouldLimit(t *testing.T) {
	config := &StoreConfig{
		MaxRequests:    10,
		LimitInSeconds: 3,
		BlockInSeconds: 60,
	}
	store := NewInMemoryStore(config)

	// Test when hit count is below the limit
	store.hitCount = 5
	if store.ShouldLimit() {
		t.Error("ShouldLimit() returned true when hit count is below the limit")
	}

	// Test when hit count is equal to the limit
	store.hitCount = 10
	if store.ShouldLimit() {
		t.Error("ShouldLimit() returned false when hit count is equal to the limit")
	}

	// Test when hit count is above the limit
	store.hitCount = 11
	if !store.ShouldLimit() {
		t.Error("ShouldLimit() returned false when hit count is above the limit")
	}
}

func TestInMemoryStore_ShouldRefresh(t *testing.T) {
	config := &StoreConfig{
		MaxRequests:    5,
		LimitInSeconds: 1,
		BlockInSeconds: 1,
	}
	store := NewInMemoryStore(config)

	store.lastHit = time.Now().Add(-100 * time.Second)
	// Test when last hit time is within the refresh interval
	if !store.ShouldRefresh() {
		t.Error("ShouldRefresh() returned false when last hit time is within the refresh interval")
	}

	store.lastHit = time.Now().Add(500 * time.Second)
	// Test when last hit time is outside the refresh interval
	if store.ShouldRefresh() {
		t.Error("ShouldRefresh() returned true when last hit time is outside the refresh interval")
	}
}

func TestInMemoryStore_Refresh(t *testing.T) {
	config := &StoreConfig{
		MaxRequests:    5,
		LimitInSeconds: 1,
		BlockInSeconds: 1,
	}
	store := NewInMemoryStore(config)

	store.hitCount = 10
	store.lastHit = time.Now().Add(-100 * time.Second)
	store.isBlocked = true
	store.Refresh()

	if store.hitCount != 1 {
		t.Error("Refresh() did not reset hit count")
	}
	if store.lastHit.Unix() != time.Now().Unix() {
		t.Error("Refresh() did not reset last hit time")
	}
	if store.isBlocked {
		t.Error("Refresh() did not reset isBlocked")
	}
}

func TestInMemoryStore_Block(t *testing.T) {
	config := &StoreConfig{
		MaxRequests:    5,
		LimitInSeconds: 1,
		BlockInSeconds: 1,
	}
	store := NewInMemoryStore(config)

	store.Block()
	if !store.isBlocked {
		t.Error("Block() did not set isBlocked")
	}
}

func TestInMemoryStore_Hit(t *testing.T) {
	config := &StoreConfig{
		MaxRequests:    5,
		LimitInSeconds: 1,
		BlockInSeconds: 1,
	}
	store := NewInMemoryStore(config)

	store.hitCount = 1
	store.lastHit = time.Now().Add(-100 * time.Second)
	store.isBlocked = false
	store.Hit()

	if store.hitCount != 2 {
		t.Error("Hit() did not increment hit count")
	}
	if store.lastHit.Unix() != time.Now().Unix() {
		t.Error("Hit() did not reset last hit time")
	}
	if store.isBlocked {
		t.Error("Hit() blocked when hit count is below the limit")
	}
	store.hitCount = 5
	store.Hit()
	if !store.isBlocked {
		t.Error("Hit() did not set isBlocked")
	}
}
