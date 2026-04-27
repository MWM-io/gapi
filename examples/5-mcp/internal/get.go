package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type getOneHandler struct {
	handler.WithMiddlewares

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

// Doc implements openapi.Documented.
// WithMCPToolName overrides the auto-generated name so the MCP tool is called "find_user"
// instead of the default "get_getUsers".
func (h getOneHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Get one").
		WithDescription("Get a user by ID").
		WithTags("Users").
		WithMCPToolName("find_user").
		WithError(404, "not_found", "user not found for id XX").
		WithResponse(User{})
	return nil
}

// Serve implements handler.Handler
func (h *getOneHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return GetByID(h.pathParameters.ID)
}
