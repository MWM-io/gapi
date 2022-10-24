package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
	"github.com/mwm-io/gapi/server/openapi"
)

// WithStatusCode is able to return its http status code.
type WithStatusCode interface {
	StatusCode() int
}

// Marshaler is able to marshal.
type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

// MarshalerFunc is function type that implements Marshaler interface.
type MarshalerFunc func(v interface{}) ([]byte, error)

// Marshal implements the Marshaler interface.
func (f MarshalerFunc) Marshal(v interface{}) ([]byte, error) {
	return f(v)
}

// ResponseWriterMiddleware is a middleware that will take the response from the next handler
// and write it into the response.
// It will choose the content type based on the request Accept header.
type ResponseWriterMiddleware struct {
	// Marshalers is the list of available Marshaler by content type.
	Marshalers map[string]Marshaler
	// DefaultContentType is the defaut content-type if the request don't have any.
	DefaultContentType string
	// Response is only use for the openAPI documentation to indicates the response type.
	Response interface{}
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

// Doc implements the openapi.OperationDescriptor interface
func (m ResponseWriterMiddleware) Doc(builder *openapi.OperationBuilder) error {
	if m.Response == nil {
		return nil
	}

	for contentType := range m.Marshalers {
		builder.WithResponse(m.Response, openapi.WithMimeType(contentType))
	}

	return builder.Error()
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
		return "", nil, errors.Wrap(errAccept, fmt.Sprintf("unknown content-type %s", accept)).WithStatus(http.StatusBadRequest)
	}

	if wantedType == "" || wantedType == "*/*" {
		wantedType = m.DefaultContentType
	}

	marshaler, ok := m.Marshalers[wantedType]
	if !ok || marshaler == nil {
		return "", nil, errors.Err(fmt.Sprintf("unsupported content-type %s", wantedType)).WithStatus(http.StatusBadRequest)
	}

	return wantedType, marshaler, nil
}

// ResponseWriter is used to store the statusCode and the content written to the http.ResponseWriter.
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
