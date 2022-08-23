package process

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// JsonBody is a pre-processor that will json.Unmarshal the request body into the Body field.
type JsonBody struct {
	Body           interface{}
	SkipValidation bool
}

// BodyValidation interface can be implemented to trigger an auto validation by JsonBody PreProcess
type BodyValidation interface {
	Validate() response.Error
}

// PreProcess implements the server.PreProcess interface
func (m JsonBody) PreProcess(handler request.Handler, r *request.WrappedRequest) (request.Handler, response.Error) {
	var buffer bytes.Buffer
	reader := io.TeeReader(r.Request.Body, &buffer)

	defer func() {
		_ = r.Request.Body.Close

		r.Request.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))
	}()

	if err := json.NewDecoder(reader).Decode(m.Body); err != nil {
		return handler, response.ErrorInternalServerErrorf("failed to parse body")
	}

	if v, ok := m.Body.(BodyValidation); !m.SkipValidation && ok {
		return handler, v.Validate()
	}

	return handler, nil
}
