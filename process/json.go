package process

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mwm-io/gapi/error"
	"github.com/mwm-io/gapi/request"
)

// JsonBody is a pre-processor that will json.Unmarshal the request body into the Body field.
type JsonBody struct {
	Body           interface{}
	SkipValidation bool
}

// BodyValidation interface can be implemented to trigger an auto validation by JsonBody PreProcess
type BodyValidation interface {
	Validate() error.Error
}

// Wrap implements the request.Middleware interface
func (m JsonBody) Wrap(h request.Handler) request.Handler {
	return request.HandlerFunc(func(r request.WrappedRequest) (interface{}, error.Error) {
		var buffer bytes.Buffer
		reader := io.TeeReader(r.Request.Body, &buffer)

		defer func() {
			_ = r.Request.Body.Close
		}()

		if err := json.NewDecoder(reader).Decode(m.Body); err != nil {
			return nil, error.Wrap(err, "failed to parse body", error.StatusCodeOpt(http.StatusBadRequest))
		}

		r.Request.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		if v, ok := m.Body.(BodyValidation); !m.SkipValidation && ok {
			if errValidate := v.Validate(); errValidate != nil {
				return nil, errValidate
			}
		}

		return h.Serve(r)
	})
}
