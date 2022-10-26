package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

// ListHandler is handler type, carrying your dependencies if you need.
type ListHandler struct {
	// Shortcut to implements the server.MiddlewareAware interface.
	server.MiddlewareHandler

	// You can also add request-scoped data.
	// That is why we need ListHandlerF HandlerFactory: a ListHandler will be created for each http request.
	queryParameters struct {
		Name string `path:"name"`
	}
}

// ListHandlerF builds the handler factory.
// You could also move this part to the main file if you prefer,
// but at MWM we prefer to keep it here to see which middlewares and options are used.
// You can also inject your dependencies here, adding parameters to the ListHandlerF function and passing them to your handler.
func ListHandlerF() server.HandlerFactory {
	return func() server.Handler {
		h := &ListHandler{}

		h.MiddlewareHandler = middleware.Core(
			middleware.WithQueryParameters(&h.queryParameters),
			middleware.WithResponseType(User{}),
		)

		return h
	}
}

func (h *ListHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return Search(h.queryParameters.Name)
}
