package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
	"github.com/mwm-io/gapi/server/openapi"

	"github.com/elnormous/contenttype"
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
	// DefaultContentType is the default content-type if the request don't have any.
	DefaultContentType string
	// ForcedContentType will always return a response serialized with this content-type.
	ForcedContentType string
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
	if m.ForcedContentType != "" {
		marshaler, ok := m.Marshalers[m.ForcedContentType]
		if !ok {
			return "", nil, errors.Err(fmt.Sprintf("no content marshaler found for content type %s", m.ForcedContentType))
		}

		return m.ForcedContentType, marshaler, nil
	}

	var availableTypes []contenttype.MediaType
	for mediaType := range m.Marshalers {
		parsedMediaType, err := contenttype.ParseMediaType(mediaType)
		if err != nil {
			return "", nil, errors.Wrap(err, fmt.Sprintf("invalid mediaType %s", mediaType)).WithStatus(http.StatusInternalServerError)
		}
		availableTypes = append(availableTypes, parsedMediaType)
	}

	accepted, _, err := contenttype.GetAcceptableMediaType(r, availableTypes)
	if err != nil {
		return "", nil, errors.Wrap(err, fmt.Sprintf("no content-type found to match the accept header %s", r.Header.Get("Accept"))).WithStatus(http.StatusUnsupportedMediaType)
	}

	marshaler, ok := m.Marshalers[accepted.String()]
	if !ok {
		return "", nil, errors.Err(fmt.Sprintf("no content marshaler found for content type %s", accepted.String()))
	}

	return accepted.String(), marshaler, nil
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
