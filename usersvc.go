package user

import "context"

type ServiceImpl struct {
	storage  Storage
	eventsvc EventService
}

// Make sure ServiceImpl implements Service
var _ Service = &ServiceImpl{}

func NewService(storage Storage, eventsvc EventService) *ServiceImpl {
	return &ServiceImpl{
		storage:  storage,
		eventsvc: eventsvc,
	}
}

// Create creates a user and publish a user.created event to our message broker.
func (s ServiceImpl) Create(ctx context.Context, u *User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Create(ctx, u)
	if err != nil {
		return err
	}
	return s.eventsvc.UserCreated(ctx, u)
}

// Update updates a user and publish a user.updated event to our message broker.
func (s ServiceImpl) Update(ctx context.Context, u *User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Update(ctx, u)
	if err != nil {
		return err
	}
	return s.eventsvc.UserUpdated(ctx, u)
}

// Delete deletes a user and publish a user.deleted event to our message broker.
func (s ServiceImpl) Delete(ctx context.Context, id string) error {
	u, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	err = s.storage.Delete(ctx, u.ID)
	if err != nil {
		return err
	}
	return s.eventsvc.UserDeleted(ctx, u)
}

// Get retrieves a user by its id.
func (s ServiceImpl) Get(ctx context.Context, id string) (*User, error) {
	return s.storage.Get(ctx, id)
}

// List retrieves a list of user using the criteria provided on opts.
func (s ServiceImpl) List(ctx context.Context, opts *ListOptions) (*ListResponse, error) {
	if opts == nil {
		opts = new(ListOptions)
	}
	return s.storage.List(ctx, opts)
}
