package process

import (
	"fmt"

	"github.com/mwm-io/gapi/error"
	"github.com/mwm-io/gapi/request"
)

// SaveRequestID /
type SaveRequestID struct{}

// PostProcess implements the server.PostProcess interface
func (m SaveRequestID) PostProcess(_ request.Handler, r *request.WrappedRequest) (request.Handler, error.Error) {

	fmt.Println(r.RequestID)

	return nil, nil
}
