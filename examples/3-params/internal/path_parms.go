package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// NewPathParamsHandler catch path params in URI, store it in h.params and reply with the params
func NewPathParamsHandler() handler.Handler {
	h := pathParamsHandler{}
	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.params},
	}

	return &h
}

type pathParamsHandler struct {
	handler.WithMiddlewares

	params struct {
		First  string `path:"first" json:"first"`
		Second string `path:"second" json:"second"`
	}
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h pathParamsHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h.params, nil
}
