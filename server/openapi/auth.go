package openapi

import (
	"net/http"
)

type Authorization interface {
	Authorize(w http.ResponseWriter, r *http.Request) (bool, error)
	Login(w http.ResponseWriter, r *http.Request) (interface{}, error)
}
