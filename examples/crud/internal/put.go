package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type PutHandler struct {
	handler.WithMiddlewares

	body           UserBody
	pathParameters struct {
		ID int `path:"id"`
	}
}

func (h PutHandler) Doc(builder *openapi.DocBuilder) error {
	builder.WithResponse(User{})
	return nil
}

func PutHandlerF() handler.Handler {
	h := &PutHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return h
}

// Serve /
func (h *PutHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user, err := GetByID(h.pathParameters.ID)
	if err != nil {
		return nil, err
	}

	user.Name = h.body.Name

	return Save(user)
}
