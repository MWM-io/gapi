package router

import (
	"os"
)

// Config describe all configurable parameters for the router server.
// This config is used by DefaultConfig
type Config struct {
	StrictSlash    bool
	Port           string
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// DefaultConfig is the default configuration used to create router & start server
//  - StrictSlash: default = true. If true it enabled automatic redirection "/path/" -> "/path"
//  - Port: default set to "PORT" env var. If this value is empty Handle use port 8080 & HandleRouterTLS use 443
//  - AllowedOrigins: default = "*". CORS config to authorize call origins
//  - AllowedMethods: default : all HTTP methods. CORS config to authorize specific HTTP methods
//  - AllowedHeaders: default = {"Content-Type", "Authorization"}. CORS config to authorize specific headers
var DefaultConfig = Config{
	StrictSlash:    true,
	Port:           os.Getenv("PORT"),
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"Content-Type", "Authorization"},
}
