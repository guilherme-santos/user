package user_test

import (
	"testing"

	"github.com/guilherme-santos/user"

	"github.com/stretchr/testify/assert"
)

func TestUserValidate(t *testing.T) {
	testcases := []struct {
		Name  string
		User  user.User
		Error string
	}{
		{
			Name:  "missing first_name",
			User:  user.User{},
			Error: "code=missing_field field=first_name message=Field was not provided",
		},
		{
			Name: "missing last_name",
			User: user.User{
				FirstName: "Guilherme",
			},
			Error: "code=missing_field field=last_name message=Field was not provided",
		},
		{
			Name: "missing password",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
			},
			Error: "code=missing_field field=password message=Field was not provided",
		},
		{
			Name: "password too weak (creating)",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "1234",
			},
			Error: "code=password_too_weak field=password message=Provided password need to be longer the 6 chars",
		},
		{
			Name: "password too weak (updating)",
			User: user.User{
				ID:        "uuid",
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "1234",
			},
			Error: "code=password_too_weak field=password message=Provided password need to be longer the 6 chars",
		},
		{
			Name: "missing email",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "123456",
			},
			Error: "code=missing_field field=email message=Field was not provided",
		},
		{
			Name: "invalid email",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "123456",
				Email:     "xguiga(at)gmail.com",
			},
			Error: "code=invalid_email field=email message=Provided email doesn't seems to be valid",
		},
		{
			Name: "missing country",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "123456",
				Email:     "xguiga@gmail.com",
			},
			Error: "code=missing_field field=country message=Field was not provided",
		},
		{
			Name: "invalid country",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "123456",
				Email:     "xguiga@gmail.com",
				Country:   "DEU",
			},
			Error: "code=invalid_country field=country message=Provided country doesn't seems to be a ISO 3166-1 alpha-2",
		},
		{
			Name: "valid for creating",
			User: user.User{
				FirstName: "Guilherme",
				LastName:  "S.",
				Password:  "123456",
				Email:     "xguiga@gmail.com",
				Country:   "DE",
			},
			Error: "",
		},
		{
			Name: "valid for updating (no password)",
			User: user.User{
				ID:        "uuid",
				FirstName: "Guilherme",
				LastName:  "S.",
				Email:     "xguiga@gmail.com",
				Country:   "DE",
			},
			Error: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.User.Validate()
			if tc.Error == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.Error)
			}
		})
	}
}
