package openapi

import (
	"net/http"
)

// Authorization is able to determine whether a request is allowed to access the documentation or not.
// It also provides a Login method to indicate how to login.
type Authorization interface {
	Authorize(w http.ResponseWriter, r *http.Request) (bool, error)
	Login(w http.ResponseWriter, r *http.Request) (interface{}, error)
}
