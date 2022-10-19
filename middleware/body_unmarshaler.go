package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
	"github.com/mwm-io/gapi/server/openapi"
)

// Unmarshaler is able to unmarshal the body into a value.
type Unmarshaler interface {
	Unmarshal(b []byte, v interface{}) error
}

// UnmarshalerFunc is function type that implements Unmarshaler interface.
type UnmarshalerFunc func(b []byte, v interface{}) error

func (f UnmarshalerFunc) Unmarshal(b []byte, v interface{}) error {
	return f(b, v)
}

// BodyUnmarshaler is a pre-processor that will Unmarshal the request body into the Body field.
type BodyUnmarshaler struct {
	Body               interface{}
	Unmarshalers       map[string]Unmarshaler
	DefaultContentType string
	SkipValidation     bool
}

// BodyValidation interface can be implemented to trigger an auto validation by BodyUnmarshaler PreProcess
type BodyValidation interface {
	Validate() error
}

// Wrap implements the request.Middleware interface
func (m BodyUnmarshaler) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Body == nil {
			return h.Serve(w, r)
		}

		unmarshaler, err := m.resolveContentType(r)
		if err != nil {
			return nil, errors.Err("unable to resolve content type", err).WithStatus(http.StatusBadRequest)
		}

		var buffer bytes.Buffer
		reader := io.TeeReader(r.Body, &buffer)

		defer func() {
			_ = r.Body.Close
		}()

		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, errors.Err("failed to read body", err).WithStatus(http.StatusBadRequest)
		}

		if errUnmarshal := unmarshaler.Unmarshal(body, m.Body); errUnmarshal != nil {
			return nil, errors.Err("failed to unmarshal body", errUnmarshal).WithStatus(http.StatusBadRequest)
		}

		r.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		if v, ok := m.Body.(BodyValidation); !m.SkipValidation && ok {
			if errValidate := v.Validate(); errValidate != nil {
				return nil, errors.Err("validation failed", errValidate).WithStatus(http.StatusUnprocessableEntity)
			}
		}

		return h.Serve(w, r)
	})
}

// Doc implements the openapi.OperationDescriptor interface
func (m BodyUnmarshaler) Doc(builder *openapi.OperationBuilder) error {
	if m.Body == nil {
		return nil
	}

	for contentType, _ := range m.Unmarshalers {
		builder.WithBody(m.Body, openapi.WithMimeType(contentType))
	}

	return builder.Error()
}

func (m BodyUnmarshaler) resolveContentType(r *http.Request) (Unmarshaler, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = m.DefaultContentType
	}

	wantedType, _, errContent := mime.ParseMediaType(contentType)
	if errContent != nil {
		return nil, errors.Err(fmt.Sprintf("unknown content-type %s", contentType), errContent)
	}

	if wantedType == "" || wantedType == "*/*" {
		wantedType = m.DefaultContentType
	}

	unmarshaler, ok := m.Unmarshalers[wantedType]
	if !ok || unmarshaler == nil {
		return nil, errors.Err(fmt.Sprintf("unsupported content-type %s", wantedType))
	}

	return unmarshaler, nil
}
