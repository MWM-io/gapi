package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

type PostHandler struct {
	server.MiddlewareHandler

	body UserBody
}

func PostHandlerF() server.HandlerFactory {
	return func() server.Handler {
		h := &PostHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithBody(&h.body),
			middleware.WithResponseType(User{}),
		)

		return h
	}
}

// Serve /
func (h *PostHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := User{
		UserBody: UserBody{
			Name: h.body.Name,
		},
	}

	return Save(user)
}
