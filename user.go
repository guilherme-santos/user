package user

import (
	"context"
	"strings"
	"time"
)

var ErrNotFound = &Err{Type: NotFound, Code: "not_found", Message: "user not found"}

//go:generate mockgen -package mock -mock_names Service=UserService -destination mock/usersvc.go github.com/guilherme-santos/user Service

// Service is an interface which implements the basic functions of the service.
type Service interface {
	// Create creates a new user.
	Create(context.Context, *User) error
	// Update updates a existing user.
	Update(context.Context, *User) error
	// Delete deletes a existing user.
	Delete(_ context.Context, id string) error
	// Get retrieves and user by id.
	Get(_ context.Context, id string) (*User, error)
	// List retrieves a list of users using the criterias on ListOptions.
	List(context.Context, *ListOptions) (*ListResponse, error)
}

//go:generate mockgen -package mock -mock_names Storage=UserStorage -destination mock/userstorage.go github.com/guilherme-santos/user Storage

// Storage is an interface which implements the storage for the service.
// In this case all methods are the same provided by Service interface.
type Storage Service

//go:generate mockgen -package mock -mock_names EventService=EventService -destination mock/eventsvc.go github.com/guilherme-santos/user EventService

// EventService is an interface which implements publishing events for outside of the
// service.
type EventService interface {
	UserCreated(context.Context, *User) error
	UserUpdated(context.Context, *User) error
	UserDeleted(context.Context, *User) error
}

// User is the representation of a user in the context of faceit
type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	// Password keeps the plain password when creating or updating a user.
	// Important: It will never be returned to the clients.
	Password  string     `json:"password,omitempty"`
	Email     string     `json:"email"`
	Country   string     `json:"country"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	RemovedAt *time.Time `json:"removed_at,omitempty"`
}

func (u User) Validate() error {
	if u.FirstName == "" {
		return NewMissingFieldError("first_name")
	}
	if u.LastName == "" {
		return NewMissingFieldError("last_name")
	}
	// No id, users is been creating, so password is required.
	if u.Password == "" && u.ID == "" {
		return NewMissingFieldError("password")
	}
	// Only check length if password is provided (creating or updating)
	if u.Password != "" && len(u.Password) < 6 {
		return &FieldError{
			Err: Error{
				Code:    "password_too_weak",
				Message: "Provided password need to be longer the 6 chars",
			},
			Field: "password",
		}
	}
	if u.Email == "" {
		return NewMissingFieldError("email")
	}
	if !strings.Contains(u.Email, "@") {
		return &FieldError{
			Err: Error{
				Code:    "invalid_email",
				Message: "Provided email doesn't seems to be valid",
			},
			Field: "email",
		}
	}
	if u.Country == "" {
		return NewMissingFieldError("country")
	}
	if len(u.Country) != 2 {
		return &FieldError{
			Err: Error{
				Code:    "invalid_country",
				Message: "Provided country doesn't seems to be a ISO 3166-1 alpha-2",
			},
			Field: "country",
		}
	}
	return nil
}

// ListOptions contains filtering, sorting and other fields to filter the list of user.
type ListOptions struct {
	Country string
	Sort    string
	PerPage int64
	Page    int
	Cursor  string
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		PerPage: 10,
	}
}

// ListResponse contains the list of users returned by List method.
type ListResponse struct {
	Total      int64   `json:"total"`
	PerPage    int64   `json:"per_page"`
	Users      []*User `json:"users"`
	NextCursor string  `json:"next_cursor"`
}
