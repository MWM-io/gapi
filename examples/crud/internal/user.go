package internal

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mwm-io/gapi/errors"
)

type User struct {
	ID int `json:"id"`
	UserBody
}

type UserBody struct {
	Name string `json:"name"`
}

func (u UserBody) Validate() error {
	if u.Name == "" {
		return errors.Err("name required").WithStatus(http.StatusPreconditionFailed)
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

func GetByID(id int) (User, error) {
	usersMu.RLock()
	defer usersMu.RUnlock()

	for _, u := range users {
		if u.ID == id {
			return u, nil
		}
	}

	return User{}, errors.Err(fmt.Sprintf("user not found for id %d", id)).
		WithKind("not_found").
		WithStatus(http.StatusNotFound)
}

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

	return errors.Err(fmt.Sprintf("user not found for id %d", id)).
		WithKind("not_found").
		WithStatus(http.StatusNotFound)
}
