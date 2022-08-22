package response

import (
	"net/http"
)

// Error TODO
type Error interface {
	Message() string
	StatusCode() int
	Unwrap() error
	WriteResponse(r http.ResponseWriter, contentType string)
}
