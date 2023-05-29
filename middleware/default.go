package middleware

import (
	"github.com/mwm-io/gapi/handler"
)

// Defaults includes all default middlewares.
// You can update this list if you want to change middleware configs for all you handlers.
var Defaults = []handler.Middleware{
	MakeJSONResponseWriter(),
	Log{},
	Recover{},
}
