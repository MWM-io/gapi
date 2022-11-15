package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// searchHandler is handler type, carrying your dependencies if you need.
type searchHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need SearchHandler Factory: a searchHandler will be created for each http request.
	queryParameters struct {
		Name string `query:"name"`
	}
}

// Doc Implements openapi.Documented: use it to set additional info about you handler
func (h *searchHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Search users").
		WithDescription("Search users by name. If query param name is empty, all users are returned").
		WithTags("Getter").
		WithResponse([]User{})
	return nil
}

// SearchHandler : Search users by name. If query param name is empty, all users are returned
func SearchHandler() handler.Handler {
	h := searchHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.QueryParameters{Parameters: &h.queryParameters},
	}

	return &h
}

// Serve Implements handler.Handler
func (h *searchHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return Search(h.queryParameters.Name)
}
