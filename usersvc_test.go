package user_test

import (
	"context"
	"testing"

	"github.com/guilherme-santos/user"
	"github.com/guilherme-santos/user/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	u := newUser()

	// Creates the user in the storage
	storage := mock.NewUserStorage(ctrl)
	storage.EXPECT().Create(gomock.Any(), u).Return(nil)

	// Publish a user.created event
	eventsvc := mock.NewEventService(ctrl)
	eventsvc.EXPECT().UserCreated(gomock.Any(), u).Return(nil)

	svc := user.NewService(storage, eventsvc)
	err := svc.Create(ctx, u)
	assert.NoError(t, err)
}

func TestUserServiceUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	u := newUser()

	// Updates the user in the storage
	storage := mock.NewUserStorage(ctrl)
	storage.EXPECT().Update(gomock.Any(), u).Return(nil)

	// Publish a user.updated event
	eventsvc := mock.NewEventService(ctrl)
	eventsvc.EXPECT().UserUpdated(gomock.Any(), u).Return(nil)

	svc := user.NewService(storage, eventsvc)
	err := svc.Update(ctx, u)
	assert.NoError(t, err)
}

func TestUserServiceDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	u := newUser()
	u.ID = "uuid"

	// Get and deletes the user in the storage
	storage := mock.NewUserStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), u.ID).Return(u, nil)
	storage.EXPECT().Delete(gomock.Any(), u.ID).Return(nil)

	// Publish a user.deleted event
	eventsvc := mock.NewEventService(ctrl)
	eventsvc.EXPECT().UserDeleted(gomock.Any(), u).Return(nil)

	svc := user.NewService(storage, eventsvc)
	err := svc.Delete(ctx, u.ID)
	assert.NoError(t, err)
}

func TestUserServiceGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	u := newUser()
	u.ID = "uuid"

	// Get and deletes the user in the storage
	storage := mock.NewUserStorage(ctrl)
	storage.EXPECT().Get(gomock.Any(), u.ID).Return(u, nil)

	eventsvc := mock.NewEventService(ctrl)

	svc := user.NewService(storage, eventsvc)
	uu, err := svc.Get(ctx, u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u, uu)
}

func TestUserServiceList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	u1 := newUser()
	u1.ID = "uuid-1"
	u2 := newUser()
	u2.ID = "uuid-2"
	u3 := newUser()
	u3.ID = "uuid-3"

	// Get and deletes the user in the storage
	storage := mock.NewUserStorage(ctrl)
	storage.EXPECT().
		List(gomock.Any(), &user.ListOptions{PerPage: 10}).
		Return(&user.ListResponse{
			Total:   2,
			PerPage: 10,
			Users:   []*user.User{u1, u2, u3},
		}, nil)

	eventsvc := mock.NewEventService(ctrl)

	svc := user.NewService(storage, eventsvc)
	list, err := svc.List(ctx, nil)
	assert.NoError(t, err)
	if assert.Len(t, list.Users, 3) {
		assert.Equal(t, u1, list.Users[0])
		assert.Equal(t, u2, list.Users[1])
		assert.Equal(t, u3, list.Users[2])
	}
}

func newUser() *user.User {
	return &user.User{
		FirstName: "Guilherme",
		LastName:  "S.",
		Password:  "123456",
		Email:     "xguiga@gmail.com",
		Country:   "DE",
	}
}
