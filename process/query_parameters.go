package process

import (
	"github.com/gorilla/schema"

	"github.com/mwm-io/gapi/error"

	"github.com/mwm-io/gapi/request"
)

var decoder = schema.NewDecoder()

// QueryParameters is a pre-processor that will set the request parameters into the Parameters field.
type QueryParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m QueryParameters) Wrap(h request.Handler) request.Handler {
	return request.HandlerFunc(func(r request.WrappedRequest) (interface{}, error.Error) {
		decoder.IgnoreUnknownKeys(true)
		decoder.SetAliasTag("query")
		err := decoder.Decode(m.Parameters, r.Request.URL.Query())
		if err != nil {
			return nil, error.Wrap(err, "decoder.Decode() failed")
		}

		return h.Serve(r)
	})
}
