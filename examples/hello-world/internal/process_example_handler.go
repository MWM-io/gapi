package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/request"
)

// ProcessBody : request body for ProcessHandler
type ProcessBody struct {
	Name string `json:"name"`
}

// ProcessResponse : result for ProcessHandler
type ProcessResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ProcessHandler /
type ProcessHandler struct {
	request.MiddlewareHandler

	body           ProcessBody
	pathParameters struct {
		ID string `path:"id"`
	}
	queryParameters struct {
		Age int `query:"age"`
	}
}

// ProcessHandlerF /
func ProcessHandlerF() request.HandlerFactory {
	return func() request.Handler {
		h := &ProcessHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithPathParameters(&h.pathParameters),
			middleware.WithQueryParameters(&h.queryParameters),
			middleware.WithBody(&h.body),
		)

		return h
	}
}

// Serve /
func (h *ProcessHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return ProcessResponse{
		ID:   h.pathParameters.ID,
		Name: h.body.Name,
		Age:  h.queryParameters.Age,
	}, nil
}
