package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
)

// JsonBody is a pre-processor that will json.Unmarshal the request body into the Body field.
type JsonBody struct {
	Body           interface{}
	SkipValidation bool
}

// BodyValidation interface can be implemented to trigger an auto validation by JsonBody PreProcess
type BodyValidation interface {
	Validate() error
}

// Wrap implements the request.Middleware interface
func (m JsonBody) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Body == nil {
			return h.Serve(w, r)
		}

		var buffer bytes.Buffer
		reader := io.TeeReader(r.Body, &buffer)

		defer func() {
			_ = r.Body.Close
		}()

		if err := json.NewDecoder(reader).Decode(m.Body); err != nil {
			return nil, errors.Wrap(err, "failed to parse body", errors.StatusCodeOpt(http.StatusBadRequest))
		}

		r.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		if v, ok := m.Body.(BodyValidation); !m.SkipValidation && ok {
			if errValidate := v.Validate(); errValidate != nil {
				return nil, errValidate
			}
		}

		return h.Serve(w, r)
	})
}
