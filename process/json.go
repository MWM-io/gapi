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

// PreProcess implements the server.PreProcess interface
func (m JsonBody) PreProcess(handler request.Handler, r *request.WrappedRequest) (request.Handler, error.Error) {
	var buffer bytes.Buffer
	reader := io.TeeReader(r.Request.Body, &buffer)

	defer func() {
		_ = r.Request.Body.Close

		r.Request.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))
	}()

	if err := json.NewDecoder(reader).Decode(m.Body); err != nil {
		return handler, error.Wrap(err, "failed to parse body", error.StatusCodeOpt(http.StatusBadRequest))
	}

	if v, ok := m.Body.(BodyValidation); !m.SkipValidation && ok {
		return handler, v.Validate()
	}

	return handler, nil
}
