package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// deleteHandler is handler type, carrying your dependencies if you need.
type deleteHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need DeleteHandler Factory: a deleteHandler will be created for each http request.
	pathParameters struct {
		ID int `path:"id"`
	}
}

// Doc Implements openapi.Documented : use it to set additional info about you handler
func (h deleteHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Delete").
		WithDescription("Delete a user by ID").
		WithError(404, "not_found", "user not found for id XX")
	return nil
}

// DeleteHandler : Delete a user by ID
func DeleteHandler() handler.Handler {
	h := deleteHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return &h
}

// Serve Implements handler.Handler
func (h *deleteHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, Delete(h.pathParameters.ID)
}
