package process

import (
	"github.com/gorilla/schema"
	"github.com/mwm-io/gapi/response"

	"github.com/mwm-io/gapi/request"
)

var decoder = schema.NewDecoder()

// QueryParameters is a pre-processor that will set the request parameters into the Parameters field.
type QueryParameters struct {
	Parameters interface{}
}

// PreProcess implements the server.PreProcess interface
func (m QueryParameters) PreProcess(handler request.Handler, r *request.WrappedRequest) (request.Handler, response.Error) {
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("query")
	err := decoder.Decode(m.Parameters, r.Request.URL.Query())
	if err != nil {
		return nil, response.ErrorInternalServerError(err)
	}

	return handler, nil
}
