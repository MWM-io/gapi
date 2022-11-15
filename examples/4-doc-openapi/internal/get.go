package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// getOneHandler is handler type, carrying your dependencies if you need.
type getOneHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need GetOneHandler Factory: a getOneHandler will be created for each http request.
	pathParameters struct {
		ID int `path:"id"`
	}
}

// GetOneHandler : Get a user by ID
func GetOneHandler() handler.Handler {
	h := getOneHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return &h
}

// Doc Implements openapi.Documented : use it to set additional info about you handler
func (h getOneHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Get one").
		WithDescription("Get a user by ID").
		WithTags("Getter").
		WithError(404, "not_found", "user not found for id XX").
		WithResponse(User{})
	return nil
}

// Serve Implements handler.Handler
func (h *getOneHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return GetByID(h.pathParameters.ID)
}
