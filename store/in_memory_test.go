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

// Add more test cases for the other methods in InMemoryStore
