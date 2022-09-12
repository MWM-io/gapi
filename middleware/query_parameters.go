package middleware

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
)

var decoder = schema.NewDecoder()

// QueryParameters is a pre-processor that will set the request parameters into the Parameters field.
type QueryParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m QueryParameters) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
