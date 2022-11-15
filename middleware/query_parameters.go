package middleware

import (
	"net/http"

	"github.com/gorilla/schema"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
)

var decoder = schema.NewDecoder()

// QueryParameters is a middleware that will set the request query parameters into the Parameters field.
type QueryParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m QueryParameters) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Parameters == nil {
			return h.Serve(w, r)
		}

		decoder.IgnoreUnknownKeys(true)
		decoder.SetAliasTag("query")
		err := decoder.Decode(m.Parameters, r.URL.Query())
		if err != nil {
			return nil, errors.Wrap(err).WithMessage("decoder.Decode() failed")
		}

		return h.Serve(w, r)
	})
}

// Doc implements the openapi.Documented interface
func (m QueryParameters) Doc(builder *openapi.DocBuilder) error {
	if m.Parameters == nil {
		return nil
	}

	return builder.WithParams(m.Parameters).Error()
}
