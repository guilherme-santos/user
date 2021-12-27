package stub

import (
	"context"

	"github.com/guilherme-santos/user"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

func (s EventService) UserCreated(ctx context.Context, u *user.User) error {
	return nil
}

func (s EventService) UserUpdated(ctx context.Context, u *user.User) error {
	return nil
}

func (s EventService) UserDeleted(ctx context.Context, u *user.User) error {
	return nil
}
