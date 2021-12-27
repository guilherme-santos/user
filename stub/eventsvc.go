package stub

import (
	"context"

	"github.com/guilherme-santos/user"

	"github.com/sirupsen/logrus"
)

type EventService struct{}

func NewEventService() *EventService {
	return &EventService{}
}

func (s EventService) UserCreated(ctx context.Context, u *user.User) error {
	s.log(ctx, "user.created", u)
	// TODO: publish event
	return nil
}

func (s EventService) UserUpdated(ctx context.Context, u *user.User) error {
	s.log(ctx, "user.updated", u)
	// TODO: publish event
	return nil
}

func (s EventService) UserDeleted(ctx context.Context, u *user.User) error {
	s.log(ctx, "user.deleted", u)
	// TODO: publish event
	return nil
}

func (s EventService) log(ctx context.Context, event string, u *user.User) {
	log := user.Logger(ctx)
	log.
		WithFields(logrus.Fields{
			"user_id": u.ID,
			"event":   event,
		}).
		Debug("publishing event")
}
