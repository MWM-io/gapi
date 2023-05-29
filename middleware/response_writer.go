package middleware

import (
	"fmt"
	"io"
	"net/http"

	"github.com/elnormous/contenttype"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/response"
)

// WithStatusCode is able to return its http status code.
type WithStatusCode interface {
	StatusCode() int
}

// ResponseWriter is a middleware that will take the response from the next handler
// and write it into the response.
// It will choose the content type based on the request Accept header.
type ResponseWriter struct {
	// Encoders is the list of available Marshaller by content type.
	Encoders map[string]Encoder
	// DefaultContentType is the default content-type if the request don't have any.
	DefaultContentType string
	// ForcedContentType will always return a response serialized with this content-type.
	ForcedContentType string
	// StatusCode stores response status code.
	StatusCode int
}

// MakeResponseWriter return an initialized ResponseWriter with all supported encoders (see EncoderByContentType)
// and ResponseWriter.DefaultContentType set to application/json
func MakeResponseWriter() ResponseWriter {
	return ResponseWriter{
		Encoders:           EncoderByContentType,
		DefaultContentType: "application/json",
		ForcedContentType:  "application/json", // TODO : Remove later
		StatusCode:         http.StatusOK,
	}
}

// MakeJSONResponseWriter return an initialized ResponseWriter with all supported encoders (see EncoderByContentType)
// and ResponseWriter.ForcedContentType set to application/json
func MakeJSONResponseWriter() ResponseWriter {
	return ResponseWriter{
		Encoders:          EncoderByContentType,
		ForcedContentType: "application/json",
		StatusCode:        http.StatusOK,
	}
}

// SetDefaultContentType set DefaultContentType and return current instance
func (m ResponseWriter) SetDefaultContentType(contentType string) ResponseWriter {
	m.DefaultContentType = contentType
	return m
}

// SetForcedContentType set ForcedContentType and return current instance
func (m ResponseWriter) SetForcedContentType(contentType string) ResponseWriter {
	m.ForcedContentType = contentType
	return m
}

// SetEncoders set Encoders and return current instance
func (m ResponseWriter) SetEncoders(encoders map[string]Encoder) ResponseWriter {
	m.Encoders = encoders
	return m
}

// AddEncoder add an Encoder in field Encoders and return current instance
func (m ResponseWriter) AddEncoder(contentType string, encoder Encoder) ResponseWriter {
	m.Encoders[contentType] = encoder
	return m
}

// Wrap implements the request.Middleware interface
func (m ResponseWriter) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		resp, err := h.Serve(w, r)

		m.StatusCode = m.StatusCodeFromHTTPServeResult(resp, err)

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

// StatusCodeFromHTTPServeResult returns response http status code from http.Serve result
func (m ResponseWriter) StatusCodeFromHTTPServeResult(resp interface{}, err error) int {
	if err != nil {
		if errWithStatus, ok := err.(WithStatusCode); ok {
			return errWithStatus.StatusCode()
		}

		return http.StatusInternalServerError
	}

	if resp == nil {
		return http.StatusNoContent
	}

	if respWithStatus, ok := resp.(WithStatusCode); ok {
		return respWithStatus.StatusCode()
	}

	switch v := resp.(type) {
	case response.Redirect:
		return v.StatusCode
	default:
		return http.StatusOK
	}
}

func (m ResponseWriter) writeResponse(w http.ResponseWriter, r *http.Request, resp interface{}) error {
	if resp == nil {
		w.WriteHeader(m.StatusCode)
		return nil
	}

	switch v := resp.(type) {
	case response.Redirect:
		http.Redirect(w, r, v.URL, v.StatusCode)
		return nil
	case io.ReadCloser:
		defer v.Close()
		w.WriteHeader(m.StatusCode)
		_, err := io.Copy(w, v)
		return err

	case io.Reader:
		w.WriteHeader(m.StatusCode)
		_, err := io.Copy(w, v)
		return err

	case []byte:
		w.WriteHeader(m.StatusCode)
		_, err := w.Write(v)
		return err

	default:
		contentType, encoder, err := m.resolveContentType(r)
		if err != nil {
			return err
		}

		w.Header().Set("Content-Type", contentType)

		body, errMarshal := encoder.Marshal(resp)
		if errMarshal != nil {
			return errMarshal
		}

		w.WriteHeader(m.StatusCode)

		_, errW := w.Write(body)

		return errW
	}
}

func (m ResponseWriter) resolveContentType(r *http.Request) (string, Encoder, error) {
	if m.ForcedContentType != "" {
		encoder, ok := m.Encoders[m.ForcedContentType]
		if !ok {
			return "", nil, errors.Err(fmt.Sprintf("no content encoder found for content type %s", m.ForcedContentType))
		}

		return m.ForcedContentType, encoder, nil
	}

	if encoder, ok := m.Encoders[m.DefaultContentType]; ok && m.DefaultContentType != "" {
		// contenttype.GetAcceptableMediaType return a random value from availableTypes in this case
		// We prefer always return the DefaultContentType
		if accept := r.Header.Get("Accept"); accept == "" || accept == "*/*" {
			return m.DefaultContentType, encoder, nil
		}
	}

	var availableTypes []contenttype.MediaType
	for mediaType := range m.Encoders {
		parsedMediaType, err := contenttype.ParseMediaType(mediaType)
		if err != nil {
			return "", nil, errors.Wrap(err).WithMessage(fmt.Sprintf("invalid mediaType %s", mediaType)).WithStatus(http.StatusInternalServerError)
		}

		availableTypes = append(availableTypes, parsedMediaType)
	}

	accepted, _, err := contenttype.GetAcceptableMediaType(r, availableTypes)
	if err != nil {
		return "", nil, errors.Wrap(err).WithMessage(fmt.Sprintf("no content-type found to match the accept header %s", r.Header.Get("Accept"))).WithStatus(http.StatusUnsupportedMediaType)
	}

	encoder, ok := m.Encoders[accepted.String()]
	if !ok {
		return "", nil, errors.Err(fmt.Sprintf("no content encoder found for content type %s", accepted.String()))
	}

	return accepted.String(), encoder, nil
}
