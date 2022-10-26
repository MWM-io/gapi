package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

type GetHandler struct {
	server.MiddlewareHandler

	pathParameters struct {
		ID int `path:"id"`
	}
}

func GetHandlerF() server.HandlerFactory {
	return func() server.Handler {
		h := &GetHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithPathParameters(&h.pathParameters),
			middleware.WithResponseType(User{}),
		)

		return h
	}
}

// Serve /
func (h *GetHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return GetByID(h.pathParameters.ID)
}
