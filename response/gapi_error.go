package response

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
)

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
func (e GapiError) WriteResponse(r http.ResponseWriter, contentType string) {
	if e.Code != 0 {
		r.WriteHeader(e.Code)
	} else {
		r.WriteHeader(http.StatusInternalServerError)
	}

	switch contentType {
	case "application/json":
		err := json.NewEncoder(r).Encode(e)
		if err != nil {
			r.WriteHeader(http.StatusInternalServerError)
			http.Error(r, err.Error(), http.StatusInternalServerError)
		}

	case "application/xml":
		b, err := xml.MarshalIndent(e, "", "	")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warn: MarshalIndent failed %s", err.Error())
			http.Error(r, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = r.Write(b)
		if err != nil {
			r.WriteHeader(http.StatusInternalServerError)
			http.Error(r, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ErrorInternalServerError returns a new InternalServerError Error. Use origin
// optional parameter to initialize the origin of this error.
func ErrorInternalServerError(origin ...error) GapiError {
	err := GapiError{
		Code: http.StatusInternalServerError,
		Msg:  "Internal server error",
	}
	if len(origin) > 0 {
		err.Origin = origin[0]
	}
	return err
}

// ErrorInternalServerErrorf returns a new InternalServerError Error with customer message.
func ErrorInternalServerErrorf(format string, args ...interface{}) GapiError {
	err := ErrorInternalServerError()
	err.Msg = fmt.Sprintf(format, args...)
	return err
}

// ErrorBadRequest returns a new BadRequest Error. Use origin
// optional parameter to initialize the origin of this error.
func ErrorBadRequest(origin ...error) GapiError {
	err := GapiError{
		Code: http.StatusBadRequest,
		Msg:  "Bad request",
	}
	if len(origin) > 0 {
		err.Origin = origin[0]
	}
	return err
}

// ErrorBadRequestf returns a new BadRequest Error with customer message.
func ErrorBadRequestf(format string, args ...interface{}) GapiError {
	err := ErrorBadRequest()
	err.Msg = fmt.Sprintf(format, args...)
	return err
}
