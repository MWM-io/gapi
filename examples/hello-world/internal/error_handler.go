package internal

import (
	"fmt"
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/request"
)

// ErrorHelloWorldHandler :
type ErrorHelloWorldHandler struct {
	request.MiddlewareHandler
}

// ErrorHelloWorldHandlerF /
func ErrorHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return ErrorHelloWorldHandler{
			MiddlewareHandler: middleware.Core(),
		}
	}
}

// Serve /
func (h ErrorHelloWorldHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	err := errors.Error{
		Msg:    "Hello World",
		Code:   500,
		Origin: fmt.Errorf("something went wrong"),
	}

	return nil, err
}
