package process

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/mwm-io/gapi/errors"

	"github.com/mwm-io/gapi/request"
)

var decoder = schema.NewDecoder()

// QueryParameters is a pre-processor that will set the request parameters into the Parameters field.
type QueryParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m QueryParameters) Wrap(h request.Handler) request.Handler {
	return request.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Parameters == nil {
			return h.Serve(w, r)
		}

		decoder.IgnoreUnknownKeys(true)
		decoder.SetAliasTag("query")
		err := decoder.Decode(m.Parameters, r.URL.Query())
		if err != nil {
			return nil, errors.Wrap(err, "decoder.Decode() failed")
		}

		return h.Serve(w, r)
	})
}
