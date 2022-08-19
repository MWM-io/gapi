package request

import (
	"net/http"

	"github.com/mwm-io/gapi/metadata"
)

// WrappedRequest /
type WrappedRequest struct {
	Request   *http.Request
	Response  http.ResponseWriter
	RequestID string
}

// NewWrappedRequest return a new instance of WrappedRequest
func NewWrappedRequest(w http.ResponseWriter, r *http.Request) WrappedRequest {
	return WrappedRequest{
		Request:   r,
		Response:  w,
		RequestID: metadata.GetRequestID(r),
	}
}
