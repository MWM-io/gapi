package server

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewServer returns a new configured *http.Server, using an existing mux.Router.
// You can override the default configuration using the available Option.
func NewServer(r *mux.Router, opts ...Option) *http.Server {
	c := newOptions(opts...)

	r = r.StrictSlash(c.StrictSlash())

	return &http.Server{
		Addr: c.Addr(),
		Handler: handlers.CORS(
			handlers.AllowedOrigins(c.CORS().AllowedOrigins),
			handlers.AllowedHeaders(c.CORS().AllowedHeaders),
			handlers.AllowedMethods(c.CORS().AllowedMethods),
		)(r),
		BaseContext: func(listener net.Listener) context.Context {
			return c.Context()
		},
	}
}

// NewMux returns a new *mux.Router.
func NewMux() *mux.Router {
	return mux.NewRouter()
}

// ServeAndHandleShutdown start a *http.Server with the default configuration (overridden by the given Option)
// This function lock your program until a signal stopping your program is received. (see WithStopSignals)
func ServeAndHandleShutdown(r *mux.Router, opts ...Option) error {
	srv := NewServer(r, opts...)

	return StartProcessAndHandleStopSignals(
		func() error {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		},
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	)
}

// ServeAndHandleTlSShutdown start a *http.Server with TLS and the default configuration (overridden by the given Option)
// This function lock your program until a signal stopping your program is received. (see WithStopSignals)
func ServeAndHandleTlSShutdown(r *mux.Router, certCRT, certKey string, opts ...Option) error {
	srv := NewServer(r, opts...)

	return StartProcessAndHandleStopSignals(
		func() error {
			if err := srv.ListenAndServeTLS(certCRT, certKey); err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		},
		func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	)
}

// StartProcessAndHandleStopSignals starts the given process and listen for os stop signals to stop it,
// executing the shutdown function.
// It can also take Option to customize StopSignals, Context, and StopTimeout.
func StartProcessAndHandleStopSignals(process func() error, shutdown func(ctx context.Context) error, opts ...Option) error {
	c := newOptions(opts...)

	processErrCh := make(chan error, 1)
	go func() {
		processErrCh <- process()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, c.StopSignals()...)

	select {
	case err := <-processErrCh:
		return err
	case <-done:
		break
	}

	ctx, cancel := context.WithTimeout(c.Context(), c.StopTimeout())
	defer cancel()

	return shutdown(ctx)
}
