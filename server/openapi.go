package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
)

// AddDocHandlers will add the necessary handlers to serve a rapidoc endpoint:
// 2 endpoints to serve rapidoc.html and oauth-receiver.html from rapidoc
// and one endpoint to serve the json openapi definition of your API.
func AddDocHandlers(r *mux.Router) error {

	AddHandler(r, http.MethodGet, openapi.Config.GetDocURI(), openapi.NewRapiDocHandler())
	AddHandler(r, http.MethodGet, openapi.Config.GetAuthReceiverURI(), handler.Func(openapi.RapiDocReceiverHandler))
	AddHandler(r, http.MethodGet, openapi.Config.GetSpecOpenAPIURI(), openapi.NewSpecOpenAPIHandler(r))

	return nil
}
