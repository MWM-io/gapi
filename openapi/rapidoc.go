package openapi

import (
	"net/http"
)

// RapiDocHandler is a server.SpecOpenAPIHandler that will serve a html page with rapidoc loading the given openAPIJsonURL
type RapiDocHandler struct {
	openAPIJsonURL string
	// auth           Authorization TODO: change by a middleware
}

// NewRapiDocHandler build a new RapiDocHandler.
func NewRapiDocHandler() RapiDocHandler {
	return RapiDocHandler{
		openAPIJsonURL: Config.GetSpecOpenAPIURI(),
	}
}

// Serve implements the server.SpecOpenAPIHandler interface
func (h RapiDocHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	// if h.auth != nil {
	// 	authorized, err := h.auth.Authorize(w, r)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !authorized {
	// 		return h.auth.Login(w, r)
	// 	}
	// }

	// TODO customisation options
	_, err := w.Write([]byte(`
<!doctype html> 
<html>
<head>
  <meta charset="utf-8"> <!-- Important: rapi-doc uses utf8 characters -->
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
  <rapi-doc
    spec-url="` + h.openAPIJsonURL + `"
	bg-color="#14191f"
	text-color="#aec2e0"
	theme="dark"
	render-style="read"
	schema-style="table"
  > </rapi-doc>
</body>
</html>`))

	return nil, err
}

// RapiDocReceiverHandler is a server.SpecOpenAPIHandler that will serve the rapidoc oauth receiver
func RapiDocReceiverHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	_, err := w.Write([]byte(`
<!doctype html>
<head>
  <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>

<body>
  <oauth-receiver> </oauth-receiver>
</body>
`))

	return nil, err
}
