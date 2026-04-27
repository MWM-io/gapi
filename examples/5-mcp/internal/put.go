package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type updateHandler struct {
	handler.WithMiddlewares

	body           UserBody
	pathParameters struct {
		ID int `path:"id"`
	}
}

// UpdateHandler : Update user by ID
func UpdateHandler() handler.Handler {
	h := updateHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return &h
}

// Doc implements openapi.Documented
func (h updateHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Update user").
		WithDescription("Update user by ID and return the updated user").
		WithTags("Users").
		WithBodyExample(UserBody{Name: "John Doe"}).
		WithResponse(User{})
	return nil
}

// Serve implements handler.Handler
func (h *updateHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user, err := GetByID(h.pathParameters.ID)
	if err != nil {
		return nil, err
	}

	user.Name = h.body.Name

	return Save(user)
}
