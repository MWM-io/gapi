package request

import (
	"io"
	"mime"
	"net/http"

	"github.com/mwm-io/gapi/errors"
)

// WithStatusCode is able to return its http status code.
type WithStatusCode interface {
	StatusCode() int
}

// Marshaler is able to marshal. TODO
type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

type MarshalerFunc func(v interface{}) ([]byte, error)

func (f MarshalerFunc) Marshal(v interface{}) ([]byte, error) {
	return f(v)
}

// ResponseWriterMiddleware /
type ResponseWriterMiddleware struct {
	Marshalers         map[string]Marshaler
	DefaultContentType string
}

// Wrap implements the request.Middleware interface
func (m ResponseWriterMiddleware) Wrap(h Handler) Handler {
	return HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		wrappedW := &ResponseWriter{ResponseWriter: w}

		resp, err := h.Serve(wrappedW, r)

		m.writeStatusCode(w, resp, err)

		var errW error
		if err != nil {
			errW = m.writeResponse(w, r, err)
		} else {
			errW = m.writeResponse(w, r, resp)
		}

		if errW != nil {
			http.Error(w, errW.Error(), http.StatusInternalServerError)
		}

		return nil, errW
	})
}

func (m ResponseWriterMiddleware) writeStatusCode(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		if errWithStatus, ok := err.(WithStatusCode); ok {
			w.WriteHeader(errWithStatus.StatusCode())
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if respWithStatus, ok := err.(WithStatusCode); ok {
		w.WriteHeader(respWithStatus.StatusCode())
		return
	}
}

func (m ResponseWriterMiddleware) writeResponse(w http.ResponseWriter, r *http.Request, resp interface{}) error {
	if resp == nil {
		return nil
	}

	switch v := resp.(type) {
	case io.ReadCloser:
		defer v.Close()
		_, err := io.Copy(w, v)
		return err

	case io.Reader:
		_, err := io.Copy(w, v)
		return err

	case []byte:
		_, err := w.Write(v)
		return err

	default:
		contentType, marshaler, err := m.resolveContentType(r)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", contentType)

		body, err := marshaler.Marshal(resp)
		if err != nil {
			return err
		}

		_, errW := w.Write(body)

		return errW
	}
}

func (m ResponseWriterMiddleware) resolveContentType(r *http.Request) (string, Marshaler, error) {
	accept := r.Header.Get("Accept")
	if accept == "" {
		accept = m.DefaultContentType
	}

	wantedType, _, errAccept := mime.ParseMediaType(accept)
	if errAccept != nil {
		return "", nil, errors.Wrap(errAccept, "unknown content-type", errors.StatusCodeOpt(http.StatusBadRequest))
	}

	if wantedType == "" || wantedType == "*/*" {
		wantedType = m.DefaultContentType
	}

	marshaler, ok := m.Marshalers[wantedType]
	if !ok || marshaler == nil {
		return "", nil, errors.Wrap(errAccept, "unsupported content-type", errors.StatusCodeOpt(http.StatusBadRequest))
	}

	return wantedType, marshaler, nil
}
