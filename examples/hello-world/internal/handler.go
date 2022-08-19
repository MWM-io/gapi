package internal

import (
	"errors"

	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// HelloWorldHandler :
type HelloWorldHandler struct {
}

// HelloWorldHandlerF /
func HelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return HelloWorldHandler{}
	}
}

// Serve /
func (h HelloWorldHandler) Serve(w request.WrappedRequest) (interface{}, response.Error) {
	// test xml response
	_ = struct {
		Text    string
		XMLName struct{} `xml:"Greetings"`
	}{Text: "Hello World"}

	// test json response
	_ = struct{ Text string }{Text: "Hello World"}

	// test error response
	err := response.GapiError{
		Msg:    "internal server error",
		Code:   500,
		Origin: errors.New("something went wrong"),
	}

	return nil, err
}
