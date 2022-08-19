package request

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"mime"
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

	result, errResp := handler.Serve(wrappedRequest)
	if errResp != nil {
		errResp.WriteResponse(w)
		return
	}

	if result == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var err error
	switch v := result.(type) {
	case io.ReadCloser:
		defer v.Close()
		_, err = io.Copy(w, v)

	case io.Reader:
		_, err = io.Copy(w, v)

	case []byte:
		w.Write(v)

	default:
		mt, _, errContent := mime.ParseMediaType(wrappedRequest.ContentType.String())
		if errContent != nil {
			http.Error(w, "malformed Content-Type header", http.StatusBadRequest)
			return
		}

		switch mt {
		case "application/xml":
			x, err := xml.MarshalIndent(result, "", "	")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.Write(x)

		case "application/json":
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(result)

		default:
			http.Error(w, "unsupported Content-Type header", http.StatusHTTPVersionNotSupported)
			return
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
