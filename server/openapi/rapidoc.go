package openapi

import (
	"net/http"

	"github.com/mwm-io/gapi/server"
)

// RapiDocHandler is a server.Handler that will serve a html page with rapidoc loading the given openAPIJsonURL
type RapiDocHandler struct {
	openAPIJsonURL string
	auth           Authorization
}

func NewRapiDocHandler(openAPIJsonURL string, auth Authorization) *RapiDocHandler {
	return &RapiDocHandler{
		openAPIJsonURL: openAPIJsonURL,
		auth:           auth,
	}
}

// Serve implements the server.Handler interface
func (h *RapiDocHandler) Serve(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	if h.auth != nil {
		authorized, err := h.auth.Authorize(w, r)
		if err != nil {
			return nil, err
		}
		if !authorized {
			return h.auth.Login(w, r)
		}
	}

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

// RapiDocReceiverHandler is a server.Handler that will serve the rapidoc oauth receiver
var RapiDocReceiverHandler = server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
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
})
