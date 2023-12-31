package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockStore struct{ mock.Mock }

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (m *MockStore) ShouldLimit() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStore) ShouldRefresh() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStore) Refresh() {
	m.Called()
}

func (m *MockStore) IsBlocked() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockStore) RemainingBlockTime() uint64 {
	args := m.Called()
	return args.Get(0).(uint64)
}

func (m *MockStore) Block() {
	m.Called()
}

func (m *MockStore) Hit() {
	m.Called()
}

func (m *MockStore) LastHit() time.Time {
	args := m.Called()
	return time.Unix(0, int64(args.Get(0).(uint64)))
}

func (m *MockStore) HitCount() uint {
	args := m.Called()
	return args.Get(0).(uint)
}
