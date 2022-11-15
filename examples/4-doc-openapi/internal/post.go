package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type createHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need CreateHandler Factory: a createHandler will be created for each http request.
	body UserBody
}

// Doc Implements openapi.Documented: use it to set additional info about you handler
func (h createHandler) Doc(builder *openapi.DocBuilder) error {
	builder.WithSummary("Create user").
		WithDescription("Create a user with given name. Return created user").
		WithBodyExample(UserBody{
			Name: "John Doe",
		}).
		WithResponse(User{})
	return nil
}

// CreateHandler : Create a user with given name. Return created user
func CreateHandler() handler.Handler {
	h := createHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
	}

	return &h
}

// Serve Implements handler.Handler
func (h createHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := User{
		UserBody: UserBody{
			Name: h.body.Name,
		},
	}

	return Save(user)
}
