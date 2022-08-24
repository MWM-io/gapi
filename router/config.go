package router

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// Config describe all configurable parameters for the http server.
type Config struct {
	port        string
	cors        *CORS
	stopTimeout *time.Duration
	stopSignals []os.Signal
}

// NewConfig builds a new configuration with the given options.
func NewConfig(opts ...ServerOption) Config {
	c := Config{}
	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// CORS contains the CORS configuration for the http server.
type CORS struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Addr return the server's address when using this configuration.
// Fallback first to the PORT environment variable then to ":8080"
func (c Config) Addr() string {
	if c.port != "" {
		return fmt.Sprintf(":%s", c.port)
	}

	if port := os.Getenv("PORT"); port != "" {
		return fmt.Sprintf(":%s", port)
	}

	return ":8080"
}

// AddrHttps return the server's address when using this configuration with https.
// Fallback first to the PORT environment variable then to ":443"
func (c Config) AddrHttps() string {
	if c.port != "" {
		return fmt.Sprintf(":%s", c.port)
	}

	if port := os.Getenv("PORT"); port != "" {
		return fmt.Sprintf(":%s", port)
	}

	return ":443"
}

// CORS return the server's CORS settings when using this configuration.
// By default, authorize "*" with methods "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"
// and headers "Content-Type", "Authorization"
func (c Config) CORS() CORS {
	if c.cors != nil {
		return *c.cors
	}

	return CORS{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}

// StopTimeout return the server's timeout when shutting down.
// Defaults to 3 seconds when not set.
func (c Config) StopTimeout() time.Duration {
	if c.stopTimeout != nil {
		return *c.stopTimeout
	}

	return 3 * time.Second
}

// StopSignals return the list of os.Signal we listen to for stopping the server.
// If empty, we use the default value: os.Interrupt, syscall.SIGINT and syscall.SIGTERM
func (c Config) StopSignals() []os.Signal {
	if len(c.stopSignals) != 0 {
		return c.stopSignals
	}

	return []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}
}

// ServerOption is an option to modify the default configuration of the http server.
type ServerOption func(*Config)

// WithPort sets the server port on which to listen to.
// If the port is empty, it will first look the PORT environment variable.
// If the PORT environment variable is empty, it will take the default value. (8080 for http and 443 for https)
func WithPort(port string) ServerOption {
	return func(config *Config) {
		config.port = port
	}
}

// WithCORS sets the cors configuration.
// If you pass nil, it will use the server default value.
func WithCORS(cors *CORS) ServerOption {
	return func(config *Config) {
		config.cors = cors
	}
}

// WithStopTimeout set the timeout when the stopping the server.
// If you pass nil, it will use the server default value.
func WithStopTimeout(d *time.Duration) ServerOption {
	return func(config *Config) {
		config.stopTimeout = d
	}
}

// WithStopSignals specify on which os.Signal we must stop the server.
func WithStopSignals(signals ...os.Signal) ServerOption {
	return func(config *Config) {
		config.stopSignals = signals
	}
}
