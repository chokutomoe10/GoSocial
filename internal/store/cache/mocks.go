package cache

import (
	"context"
	"log"
	"p1/internal/store"

	"github.com/stretchr/testify/mock"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := m.Called(userID)
	log.Printf("MockCacheStore.Get called with userID: %d\n", userID)
	return nil, args.Error(1)
}

func (m *MockUserStore) Set(ctx context.Context, user *store.User) error {
	args := m.Called(user)
	log.Printf("MockCacheStore.Set called with user: %v\n", user)
	return args.Error(0)
}
