package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type createHandler struct {
	handler.WithMiddlewares

	body UserBody
}

// CreateHandler : Create a user with given name
func CreateHandler() handler.Handler {
	h := createHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
	}

	return &h
}

// Doc implements openapi.Documented
func (h createHandler) Doc(builder *openapi.DocBuilder) error {
	builder.
		WithSummary("Create user").
		WithDescription("Create a user with given name. Return created user").
		WithTags("Users").
		WithBodyExample(UserBody{Name: "John Doe"}).
		WithResponse(User{})
	return nil
}

// Serve implements handler.Handler
func (h createHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return Save(User{UserBody: UserBody{Name: h.body.Name}})
}
