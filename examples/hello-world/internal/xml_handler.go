package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/process"
	"github.com/mwm-io/gapi/request"
)

// XmlHelloWorldHandler : Xml hello returns an xml response
type XmlHelloWorldHandler struct {
	request.MiddlewareHandler
}

// XmlHelloWorldHandlerF /
func XmlHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return XmlHelloWorldHandler{
			MiddlewareHandler: process.Core(
				process.WithDefaultContentType("application/xml"),
			),
		}
	}
}

// Serve /
func (h XmlHelloWorldHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return struct {
		Text    string
		XMLName struct{} `xml:"Greetings"`
	}{Text: "Hello World"}, nil
}
