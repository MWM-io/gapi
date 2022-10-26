package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

type DeleteHandler struct {
	server.MiddlewareHandler

	pathParameters struct {
		ID int `path:"id"`
	}
}

func DeleteHandlerF() server.HandlerFactory {
	return func() server.Handler {
		h := &DeleteHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithPathParameters(&h.pathParameters),
			middleware.WithResponseType(User{}),
		)

		return h
	}
}

// Serve /
func (h *DeleteHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return nil, Delete(h.pathParameters.ID)
}
