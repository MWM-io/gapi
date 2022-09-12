package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/middleware"
	"github.com/mwm-io/gapi/server"
)

// XmlHelloWorldHandler : Xml hello returns an xml response
type XmlHelloWorldHandler struct {
	server.MiddlewareHandler
}

// XmlHelloWorldHandlerF /
func XmlHelloWorldHandlerF() server.HandlerFactory {
	return func() server.Handler {
		return XmlHelloWorldHandler{
			MiddlewareHandler: middleware.Core(
				middleware.WithDefaultContentType("application/xml"),
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
