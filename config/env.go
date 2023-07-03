package config

import (
	"os"
)

// PORT is the port used by the server to listen
// default value is 8080
var PORT string

// IS_LOCAL is a flag to indicate if the server is running locally
var IS_LOCAL = os.Getenv("IS_LOCAL") == "true"

func init() {
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}
}
