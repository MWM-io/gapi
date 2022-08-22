package request

import (
	"net/http"

	"github.com/mwm-io/gapi/response"
)

// Handler is able to respond to a http request
type Handler interface {
	Serve(WrappedRequest) (interface{}, response.Error)
}

// HandlerFactory is a function that return a new Handler
type HandlerFactory func() Handler

type httpHandler struct {
	factory HandlerFactory
}

// ServeHTTP implements the http.Handler interface
func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.factory()

	wrappedRequest := NewWrappedRequest(w, r)

	// TODO PreProcess

	w.Header().Set("Content-Type", wrappedRequest.ContentType.String())
	result, errResp := handler.Serve(wrappedRequest)
	if errResp != nil {
		errResp.WriteResponse(wrappedRequest.Response, wrappedRequest.ContentType.String())
		return
	}

	handleHandlerResponse(wrappedRequest, result)
}
