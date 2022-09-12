package middleware

import (
	"bytes"
	"io"
	"mime"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
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
func (m ResponseWriterMiddleware) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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

// ResponseWriter This ResponseWriter is used to store the statusCode
// and the content written to the http.ResponseWriter.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	content    *bytes.Buffer
}

// NewResponseWriter Create a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{
		ResponseWriter: w,
		statusCode:     0,
		content:        new(bytes.Buffer),
	}
}

// StatusCode return the statusCode of the response.
// It's 0 if it isn't set.
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Content returns the content already written to the response.
func (rw *ResponseWriter) Content() io.Reader {
	return rw.content
}

// WriteHeader Write the code in local and to the http response
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	rw.content.Write(b)
	return rw.ResponseWriter.Write(b)
}
