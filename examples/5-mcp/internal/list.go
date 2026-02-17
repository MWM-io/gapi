package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type searchHandler struct {
	handler.WithMiddlewares

	queryParameters struct {
		Name string `query:"name"`
	}
}

// SearchHandler : Search users by name
func SearchHandler() handler.Handler {
	h := searchHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.QueryParameters{Parameters: &h.queryParameters},
	}

	return &h
}

// Doc implements openapi.Documented
func (h *searchHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Search users").
		WithDescription("Search users by name. If query param name is empty, all users are returned").
		WithTags("Users").
		WithResponse([]User{})
	return nil
}

// Serve implements handler.Handler
func (h *searchHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return Search(h.queryParameters.Name)
}
