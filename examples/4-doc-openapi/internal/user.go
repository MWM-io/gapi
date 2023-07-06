package internal

import (
	"strings"
	"sync"

	"github.com/mwm-io/gapi/errors"
)

// User model
// Tag description is used by openapi spec generation to add explanations about the field
type User struct {
	ID int `json:"id" description:"Immutable & unique user identifier"`
	UserBody
}

// UserBody model used as body by CreateHandler & UpdateHandler
//
// - Tag description is used by openapi spec generation to add explanations about the field
// - Tag required is used by middleware.Body for body validation & openapi spec generation to flag the field as required
type UserBody struct {
	Name string `json:"name" description:"User name: can be updated" required:"true"`
}

// Validate is an implementation of middleware.BodyValidation. If user UserBody is given to
// middleware.Body, this function was called automatically and error handled
func (u UserBody) Validate() error {
	if u.Name == "" {
		return errors.PreconditionFailed("missing_name", "name required")
	}

	return nil
}

// Below is an in-memory implementation of the users database for this example.

var usersMu sync.RWMutex
var users = []User{
	{
		ID: 1,
		UserBody: UserBody{
			Name: "John Smith",
		},
	},
	{
		ID: 2,
		UserBody: UserBody{
			Name: "Bridget Jones",
		},
	},
	{
		ID: 3,
		UserBody: UserBody{
			Name: "Juliet Quinn",
		},
	},
	{
		ID: 4,
		UserBody: UserBody{
			Name: "Frank Davis",
		},
	},
}

// GetByID returns the user with the given ID, or an error if not found.
func GetByID(id int) (User, error) {
	usersMu.RLock()
	defer usersMu.RUnlock()

	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}

	return User{}, errors.NotFound("user_not_found", "user not found for id %d", id)
}

// Search returns the users whose name contains the given string.
func Search(name string) ([]User, error) {
	usersMu.RLock()
	defer usersMu.RUnlock()

	var results []User
	for _, u := range users {
		if strings.Contains(u.Name, name) {
			results = append(results, u)
		}
	}

	return results, nil
}

// Save saves the given user, assigning it a new ID if it is a new user.
func Save(user User) (User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if user.ID == 0 {
		maxID := 0
		for _, u := range users {
			if u.ID > maxID {
				maxID = u.ID
			}
		}

		user.ID = maxID + 1
	} else {
		for i, u := range users {
			if u.ID == user.ID {
				users[i] = user
				return user, nil
			}
		}
	}

	users = append(users, user)

	return user, nil
}

// Delete deletes the user with the given ID, or returns an error if not found.
func Delete(id int) error {
	usersMu.RLock()
	defer usersMu.RUnlock()

	for i, u := range users {
		if u.ID == id {
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			return nil
		}
	}

	return errors.NotFound("user_not_found", "user not found for id %d", id)
}
