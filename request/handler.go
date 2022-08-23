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

	if preHandler, ok := handler.(PreProcessAware); ok {
		var errPreProcess error.Error
		for _, preprocess := range preHandler.PreProcesses() {
			handler, errPreProcess = preprocess.PreProcess(handler, &wrappedRequest)
			if errPreProcess != nil {
				handleError(wrappedRequest, errPreProcess)
				return
			}
		}
	}

	result, errResp := handler.Serve(wrappedRequest)
	if errResp != nil {
		handleError(wrappedRequest, errResp)
		return
	}

	handleHandlerResponse(wrappedRequest, result)

	if preHandler, ok := handler.(PostProcessAware); ok {
		var errPreProcess error.Error
		for _, preprocess := range preHandler.PostProcesses() {
			handler, errPreProcess = preprocess.PostProcess(handler, &wrappedRequest)
			if errPreProcess != nil {
				handleError(wrappedRequest, errPreProcess)
				return
			}
		}
	}

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
