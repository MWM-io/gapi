package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/openapi"
)

// ListHandler is handler type, carrying your dependencies if you need.
type ListHandler struct {
	// Shortcut to implement the server.MiddlewareAware interface.
	handler.WithMiddlewares

	// You can also add request-scoped data.
	// That is why we need ListHandlerF Factory: a ListHandler will be created for each http request.
	queryParameters struct {
		Name string `path:"name"`
	}
}

func (h ListHandler) Doc(builder *openapi.DocBuilder) error {
	builder.WithResponse([]User{})
	return nil
}

// ListHandlerF builds the handler factory.
// You could also move this part to the main file if you prefer,
// but at MWM we prefer to keep it here to see which middlewares and options are used.
// You can also inject your dependencies here, adding parameters to the ListHandlerF function and passing them to your handler.
func ListHandlerF() handler.Handler {
	h := &ListHandler{}

	h.MiddlewareList = []handler.Middleware{
		middleware.QueryParameters{Parameters: &h.queryParameters},
	}

	return h
}

func (h *ListHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return Search(h.queryParameters.Name)
}
