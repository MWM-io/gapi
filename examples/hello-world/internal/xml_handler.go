package internal

import (
	"github.com/mwm-io/gapi/error"
	"github.com/mwm-io/gapi/request"
)

// XmlHelloWorldHandler : Xml hello returns an xml response
type XmlHelloWorldHandler struct {
}

// XmlHelloWorldHandlerF /
func XmlHelloWorldHandlerF() request.HandlerFactory {
	return func() request.Handler {
		return XmlHelloWorldHandler{}
	}
}

// Serve /
func (h XmlHelloWorldHandler) Serve(w request.WrappedRequest) (interface{}, error.Error) {
	w.ContentType = request.ApplicationXML

	return struct {
		Text    string
		XMLName struct{} `xml:"Greetings"`
	}{Text: "Hello World"}, nil
}
