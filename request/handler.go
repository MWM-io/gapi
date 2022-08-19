package request

import (
	"encoding/json"
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

	result, err := handler.Serve(wrappedRequest)
	if err != nil {
		print(err.Message())
		// TODO handle error
		return
	}

	respBytes, errMarshal := json.Marshal(result)
	if errMarshal != nil {
		print(errMarshal.Error())
	}

	// TODO handle response
	print(string(respBytes))
}
