package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

type PutHandler struct {
	server.MiddlewareHandler

	body           UserBody
	pathParameters struct {
		ID int `path:"id"`
	}
}

func PutHandlerF() server.HandlerFactory {
	return func() server.Handler {
		h := &PutHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithBody(&h.body),
			middleware.WithPathParameters(&h.pathParameters),
			middleware.WithResponseType(User{}),
		)

		return h
	}
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
