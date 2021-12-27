package user

import "context"

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

// Storage is an interface which implements the storage for the service.
// In this case all methods are the same provided by Service interface.
type Storage Service

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
	Password  string  `json:"password,omitempty"`
	Email     string  `json:"email"`
	Country   string  `json:"country"`
	CreatedAt string  `json:"created_at"`
	UpdatedBy *string `json:"updated_by"`
}

// ListOptions contains filtering, sorting and other fields to filter the list of user.
type ListOptions struct {
	Country string
	Sort    string
	PerPage int64
	Page    int
	Cursor  string
}

// ListResponse contains the list of users returned by List method.
type ListResponse struct {
	Total      int64   `json:"total"`
	PerPage    int64   `json:"per_page"`
	Users      []*User `json:"users"`
	NextCursor string  `json:"next_cursor"`
}
