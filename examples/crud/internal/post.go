package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type PostHandler struct {
	handler.WithMiddlewares

	body UserBody
}

func (h PostHandler) Doc(builder *openapi.DocBuilder) error {
	builder.WithResponse(User{})
	return nil
}

func PostHandlerF() handler.Handler {
	h := PostHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
	}

	return h
}

// Serve /
func (h PostHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := User{
		UserBody: UserBody{
			Name: h.body.Name,
		},
	}

	return Save(user)
}
