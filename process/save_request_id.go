package process

import (
	"fmt"

	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// SaveRequestID /
type SaveRequestID struct{}

// PostProcess implements the server.PostProcess interface
func (m SaveRequestID) PostProcess(_ request.Handler, r *request.WrappedRequest) (request.Handler, response.Error) {

	fmt.Println(r.RequestID)

	return nil, nil
}
