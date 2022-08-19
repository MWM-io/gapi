package response

import (
	"encoding/json"
	"net/http"
)

// Error TODO
type Error interface {
	Message() string
	StatusCode() int
	Unwrap() error
	WriteResponse(http.ResponseWriter)
}

// GapiError /
type GapiError struct {
	Msg    string `json:"message"`
	Code   int    `json:"code"`
	Origin error  `json:"-"`
}

// Message returns the Error message. If the Msg field is not filled,
// try to call the origin Error method instead.
func (e GapiError) Message() string {
	if e.Msg == "" && e.Origin != nil {
		return e.Origin.Error()
	}
	return e.Msg
}

// StatusCode /
func (e GapiError) StatusCode() int {
	return e.Code
}

// Unwrap /
func (e GapiError) Unwrap() error {
	return e.Origin
}

// WriteResponse /
func (e GapiError) WriteResponse(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(e.Code)

	err := json.NewEncoder(rw).Encode(e)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
