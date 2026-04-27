package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type deleteHandler struct {
	handler.WithMiddlewares

	pathParameters struct {
		ID int `path:"id"`
	}
}

// DeleteHandler : Delete a user by ID
func DeleteHandler() handler.Handler {
	h := deleteHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return &h
}

// Doc implements openapi.Documented
func (h deleteHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Delete user").
		WithDescription("Delete a user by ID").
		WithTags("Users").
		WithError(404, "not_found", "user not found for id XX")
	return nil
}

// Serve implements handler.Handler
func (h *deleteHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, Delete(h.pathParameters.ID)
}
