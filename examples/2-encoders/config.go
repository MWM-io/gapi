package main

import (
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

func init() {
	middleware.Defaults = []handler.Middleware{
		// We removed the default response writer to control response encoding for each handler
		middleware.Log{},
		middleware.Recover{},
	}
}
