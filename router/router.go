package router

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Create returns a mux.Router configured with strictSlash as DefaultConfig.StrictSlash by default.
//  You can set DefaultConfig.StrictSlash = false to disable automatic redirection "/path/" -> "/path"
func Create() *mux.Router {
	return mux.NewRouter().StrictSlash(DefaultConfig.StrictSlash)
}

// Handle start server
//  - Compute CORS
//  - Listen on port defined by in DefaultConfig.Port (env variable "PORT"), if empty the port 8080 is used
//  - Serve incoming request
//  This function lock your program until a SIGINT or SIGTERM is sent, you can use this behavior to detect the server shutdown
func Handle(r *mux.Router) {
	if DefaultConfig.Port == "" {
		DefaultConfig.Port = "8080"
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", DefaultConfig.Port),
		Handler: computeCors(r),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// TODO : mLogger.Infof(nil, "", fmt.Sprintf("Server Started on %s", srv.Addr))

	handleSignalsToStopServer(srv)
}

// HandleTLS start server
//  - Compute CORS
//  - Listen on port defined by in DefaultConfig.Port (env variable "PORT"), if empty the port 443 is used
//  - Serve incoming TLS request
//  This function lock your program until a SIGINT or SIGTERM is sent, you can use this behavior to detect the server shutdown
func HandleTLS(r *mux.Router, certCRT, certKey string) {
	if DefaultConfig.Port == "" {
		DefaultConfig.Port = "443"
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", DefaultConfig.Port),
		Handler: computeCors(r),
	}

	go func() {
		if err := srv.ListenAndServeTLS(certCRT, certKey); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// TODO : mLogger.Infof(nil, "", fmt.Sprintf("Server Started on %s", srv.Addr))

	handleSignalsToStopServer(srv)
}

func handleSignalsToStopServer(srv *http.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		// TODO : mLogger.Errorf(nil, "", "Server Shutdown Failed:%+v", err)
	}

	// TODO : mLogger.Infof(nil, "", "Server stopped")
}

func computeCors(r *mux.Router) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins(DefaultConfig.AllowedOrigins),
		handlers.AllowedHeaders(DefaultConfig.AllowedHeaders),
		handlers.AllowedMethods(DefaultConfig.AllowedMethods),
	)(r)
}
