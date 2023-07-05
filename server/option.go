package server

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/mwm-io/gapi/config"
)

// Option is an option to modify the default configuration of the http server.
type Option func(*serverOptions)

// WithPort sets the server port on which to listen to.
// If the port is empty, it will first look the PORT environment variable.
// If the PORT environment variable is empty, it will take the default value. (8080 for http and 443 for https)
func WithPort(port string) Option {
	return func(config *serverOptions) {
		config.port = port
	}
}

// CORS contains the CORS configuration for the http server.
type CORS struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// WithCORS sets the cors configuration.
// By default, authorize "*" with methods "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"
// and headers "Content-Type", "Authorization"
func WithCORS(cors CORS) Option {
	return func(config *serverOptions) {
		config.cors = &cors
	}
}

// WithStopTimeout set the timeout when the shutting down the server.
// By default, 3 seconds.
func WithStopTimeout(d time.Duration) Option {
	return func(config *serverOptions) {
		config.stopTimeout = &d
	}
}

// WithStopSignals specify on which os.Signal we must shut down the server.
// By default, os.Interrupt, syscall.SIGINT and syscall.SIGTERM
func WithStopSignals(signals ...os.Signal) Option {
	return func(config *serverOptions) {
		config.stopSignals = signals
	}
}

// WithStrictSlash specify if we use the strictSlash configuration or not.
// When true, if the route path is "/path/", accessing "/path" will perform a redirect
// to the former and vice versa. In other words, your application will always
// see the path as specified in the route.
func WithStrictSlash(strictSlash bool) Option {
	return func(config *serverOptions) {
		config.withoutStrictSlash = !strictSlash
	}
}

// WithContext specify a parent context for the *http.Server.
// It will be passed to every request.
// By default, creates a new context.
func WithContext(ctx context.Context) Option {
	return func(config *serverOptions) {
		config.context = ctx
	}
}

type serverOptions struct {
	withoutStrictSlash bool
	port               string
	cors               *CORS
	stopTimeout        *time.Duration
	stopSignals        []os.Signal
	context            context.Context
}

func newOptions(opts ...Option) serverOptions {
	c := serverOptions{}
	for _, opt := range opts {
		opt(&c)
	}

	return c
}

// Addr returns the server address to listen to.
func (c serverOptions) Addr() string {
	if c.port != "" {
		return fmt.Sprintf(":%s", c.port)
	}

	return fmt.Sprintf(":%s", config.PORT)
}

// AddrHttps returns the server address to listen to.
func (c serverOptions) AddrHttps() string {
	if c.port != "" {
		return fmt.Sprintf(":%s", c.port)
	}

	return fmt.Sprintf(":%s", config.PORT)
}

// CORS returns the CORS configuration.
func (c serverOptions) CORS() CORS {
	if c.cors != nil {
		return *c.cors
	}

	return CORS{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}
}

// StopTimeout returns the time to wait for before shutting down the server.
func (c serverOptions) StopTimeout() time.Duration {
	if c.stopTimeout != nil {
		return *c.stopTimeout
	}

	return 3 * time.Second
}

// StopSignals returns the signals to listen to for shutting down the server.
func (c serverOptions) StopSignals() []os.Signal {
	if len(c.stopSignals) != 0 {
		return c.stopSignals
	}

	return []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM}
}

// StrictSlash returns the strictSlash configuration.
func (c serverOptions) StrictSlash() bool {
	return !c.withoutStrictSlash
}

// Context returns the parent context for the *http.Server.
func (c serverOptions) Context() context.Context {
	if c.context != nil {
		return c.context
	}

	return context.Background()
}
