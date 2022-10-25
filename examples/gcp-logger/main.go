package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
	"github.com/mwm-io/gapi/middleware"

	"github.com/mwm-io/gapi/server/openapi"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/log/cloud_logging"
	"github.com/mwm-io/gapi/server"
)

func main() {
	ctx := context.Background()
	clientLogger := setupLog(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	defer clientLogger.Close()

	r := server.NewMux()

	server.AddHandler(r, "GET", "/hello", server.HandlerF(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		return "hello", nil
	}, middleware.Core().Middlewares()...))

	err := openapi.AddRapidocHandlers(r, openapi.Config{})
	if err != nil {
		log.Printf("error while adding rapidoc %+v\n", err)
	}

	gLog.Info("Starting http server")

	if err := server.ServeAndHandleShutdown(r, server.WithContext(ctx)); err != nil {
		gLog.LogAny(err)
	}

	gLog.Info("Server stopped")
}

func setupLog(ctx context.Context, projectID string) *logging.Client {
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalln(err)
	}

	// This logger will both write on your console and on google cloud logging.
	log := gLog.NewDefaultLogger(
		gLog.NewMultiWriter(
			gLog.NewFilterWriter(gLog.InfoSeverity, cloud_logging.NewWriter(client.Logger("application"), projectID)),
			gLog.NewWriter(gLog.EntryMarshalerFunc(func(entry gLog.Entry) []byte {
				return []byte(fmt.Sprintf("%-9s %s %s| \n", entry.Severity, entry.Timestamp.Format("15:05:05.999999999"), entry.Message))
			}), os.Stdout),
		),
	)

	log = log.WithLabels(map[string]string{"PROJECT_ID": projectID})
	gLog.SetGlobalLogger(log)

	return client
}
