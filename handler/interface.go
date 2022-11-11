package handler

import (
	"net/http"
)

// Handler is able to respond to a http request and return the response it wants to write and an error.
// You need to use a middleware to write your response if you just return your response and not write it.
// (see github.com/mwm-io/gapi/middleware.ResponseWriter)
type Handler interface {
	Serve(http.ResponseWriter, *http.Request) (interface{}, error)
}
