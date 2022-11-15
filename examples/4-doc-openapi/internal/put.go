package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// updateHandler is handler type, carrying your dependencies if you need.
type updateHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need UpdateHandler Factory: a updateHandler will be created for each http request.
	body           UserBody
	pathParameters struct {
		ID int `path:"id"`
	}
}

// Doc Implements openapi.Documented: use it to set additional info about you handler
func (h updateHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Update user").
		WithDescription("Update user by ID and return the updated user. If user doesn't exist a new entry is created.").
		WithBodyExample(UserBody{
			Name: "John Doe",
		}).WithResponse(User{})
	return nil
}

// UpdateHandler : Update user by ID and return the updated user. If user doesn't exist a new entry is created.
func UpdateHandler() handler.Handler {
	h := updateHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return &h
}

// Serve Implements handler.Handler
func (h *updateHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user, err := GetByID(h.pathParameters.ID)
	if err != nil {
		return nil, err
	}

	user.Name = h.body.Name

	return Save(user)
}
