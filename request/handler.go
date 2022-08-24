package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"

	"github.com/mwm-io/gapi/error"
)

// Handler is able to respond to a http request
type Handler interface {
	Serve(WrappedRequest) (interface{}, error.Error)
}

type HandlerFunc func(WrappedRequest) (interface{}, error.Error)

func (h HandlerFunc) Serve(request WrappedRequest) (interface{}, error.Error) {
	return h(request)
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
	w.Header().Set("Content-Type", wrappedRequest.ContentType.String())

	if middlewareHandler, ok := handler.(MiddlewareAware); ok {
		middlewares := middlewareHandler.Middlewares()
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i].Wrap(handler)
		}
	}

	result, errResp := handler.Serve(wrappedRequest)
	if errResp != nil {
		handleError(wrappedRequest, errResp)
		return
	}

	handleHandlerResponse(wrappedRequest, result)
}

func handleError(wr WrappedRequest, errE error.Error) {
	if errE.StatusCode() != 0 {
		wr.Response.WriteHeader(errE.StatusCode())
	} else {
		wr.Response.WriteHeader(http.StatusInternalServerError)
	}

	switch wr.ContentType {
	case "application/json":
		err := json.NewEncoder(wr.Response).Encode(errE)
		if err != nil {
			http.Error(wr.Response, err.Error(), http.StatusInternalServerError)
		}

	case "application/xml":
		b, err := xml.MarshalIndent(errE, "", "	")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warn: MarshalIndent failed %s", err.Error())
			http.Error(wr.Response, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = wr.Response.Write(b)
		if err != nil {
			http.Error(wr.Response, err.Error(), http.StatusInternalServerError)
		}
	}
}
