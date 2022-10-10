package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// NewMux returns a new *mux.Router.
func NewMux() *mux.Router {
	return mux.NewRouter()
}

// AddHandler add a new handler factory to mux router for given method and path
func AddHandler(router *mux.Router, method, path string, f Handler) {
	router.Methods(method).
		Path(path).
		Handler(HttpHandler{f})
}

// ServeAndHandleShutdown start a *http.Server with the default configuration (overridden by the given options)
// It will add a CORS middleware to the *mux.Router
// This function lock your program until a signal stopping your program is received. (see WithStopSignals)
func ServeAndHandleShutdown(r *mux.Router, opts ...ServerOption) error {
	srv := NewServer(r, opts...)

	log.Printf("Server Started on %s\n", srv.Addr)

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

// ServeAndHandleTlsShutdown start a *http.Server with the default configuration (overridden by the given options) and TLS
// It will add a CORS middleware to the *mux.Router
// This function lock your program until a signal stopping your program is received. (see WithStopSignals)
func ServeAndHandleTlsShutdown(r *mux.Router, certCRT, certKey string, opts ...ServerOption) error {
	srv := NewServer(r, opts...)

	log.Printf("Server Started on %s\n", srv.Addr)

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

// NewServer returns a new configured *http.Server.
// You can override the default configuration using the available ServerOption.
func NewServer(r *mux.Router, opts ...ServerOption) *http.Server {
	c := newConfig(opts...)

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

// StartProcessAndHandleStopSignals starts the given server and handles the os stop signal to stop it,
// executing the shutdown function.
func StartProcessAndHandleStopSignals(process func() error, shutdown func(ctx context.Context) error, opts ...ServerOption) error {
	c := newConfig(opts...)

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