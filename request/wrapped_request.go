package request

import (
	"net/http"
)

// WrappedRequest /
type WrappedRequest struct {
	Request     *http.Request
	Response    http.ResponseWriter
	ContentType ContentType
}

// ContentType /
type ContentType string

// String implements the Stringer interface
func (c ContentType) String() string {
	return string(c)
}

var (
	// ApplicationJSON /
	ApplicationJSON ContentType = "application/json"
	// ApplicationXML /
	ApplicationXML ContentType = "application/xml"
)

// NewWrappedRequest return a new instance of WrappedRequest
func NewWrappedRequest(w http.ResponseWriter, r *http.Request) WrappedRequest {
	return WrappedRequest{
		Request:     r,
		Response:    w,
		ContentType: DefaultConfig.ContentType,
	}
}
