package mysql

import (
	"context"

	"github.com/guilherme-santos/user"
)

type UserStorage struct {
}

func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

func (s UserStorage) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (s UserStorage) Update(ctx context.Context, u *user.User) error {
	return nil
}

func (s UserStorage) Delete(ctx context.Context, id string) error {
	return nil
}

func (s UserStorage) Get(ctx context.Context, id string) (*user.User, error) {
	return &user.User{
		ID:        id,
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}, nil
}

func (s UserStorage) List(ctx context.Context, opts *user.ListOptions) (*user.ListResponse, error) {
	u1 := &user.User{
		ID:        "uuid-1",
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}
	u2 := &user.User{
		ID:        "uuid-2",
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}
	u3 := &user.User{
		ID:        "uuid-3",
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}

	return &user.ListResponse{
		Total:   2,
		PerPage: 10,
		Users: []*user.User{
			u1, u2, u3,
		},
	}, nil
}
