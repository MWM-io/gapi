package hello

import (
	"net/http"

	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// MakeXMLResponseHandler return a Hello world XML
func MakeXMLResponseHandler() handler.Handler {
	return xmlResponseHandler{
		WithMiddlewares: handler.WithMiddlewares{
			MiddlewareList: []handler.Middleware{
				middleware.MakeResponseWriter().
					SetForcedContentType("application/xml"),
			},
		},
	}
}

type xmlResponseHandler struct {
	handler.WithMiddlewares
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h xmlResponseHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return "Hello world", nil
}
