package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

type GetHandler struct {
	handler.WithMiddlewares

	pathParameters struct {
		ID int `path:"id"`
	}
}

func (h GetHandler) Doc(builder *openapi.DocBuilder) error {
	builder.WithResponse(User{})
	return nil
}

func GetHandlerF() handler.Handler {
	h := &GetHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return h
}

// Serve /
func (h *GetHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return GetByID(h.pathParameters.ID)
}
