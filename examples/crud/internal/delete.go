package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

type DeleteHandler struct {
	handler.WithMiddlewares

	pathParameters struct {
		ID int `path:"id"`
	}
}

func DeleteHandlerF() handler.Handler {
	h := &DeleteHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.PathParameters{Parameters: &h.pathParameters},
	}

	return h
}

// Serve /
func (h *DeleteHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, Delete(h.pathParameters.ID)
}
