package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// NewQueryParamsHandler catch the registered query params, store it in h.params and reply with the params
func NewQueryParamsHandler() handler.Handler {
	h := queryParamsHandler{}
	h.MiddlewareList = []handler.Middleware{
		middleware.QueryParameters{Parameters: &h.params},
	}

	return &h
}

type queryParamsHandler struct {
	handler.WithMiddlewares

	params struct {
		First  string `query:"first" json:"first"`
		Second string `query:"second" json:"second"`
	}
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h queryParamsHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return h.params, nil
}
