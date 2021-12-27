package stub

import (
	"context"

	"github.com/guilherme-santos/user"
)

type UserStorageCache struct {
	storage user.Storage
}

func NewUserStorageCache(storage user.Storage) *UserStorageCache {
	return &UserStorageCache{
		storage: storage,
	}
}

func (c UserStorageCache) Create(ctx context.Context, u *user.User) error {
	// TODO: invalidate the cache
	return c.storage.Create(ctx, u)
}

func (c UserStorageCache) Update(ctx context.Context, u *user.User) error {
	// TODO: invalidate the cache
	return c.storage.Update(ctx, u)
}

func (c UserStorageCache) Delete(ctx context.Context, id string) error {
	// TODO: invalidate the cache
	return c.storage.Delete(ctx, id)
}

func (c UserStorageCache) Get(ctx context.Context, id string) (*user.User, error) {
	// TODO: Save user in the cache and also etag (if implemented)
	return c.storage.Get(ctx, id)
}

func (c UserStorageCache) List(ctx context.Context, opts *user.ListOptions) (*user.ListResponse, error) {
	return c.storage.List(ctx, opts)
}
