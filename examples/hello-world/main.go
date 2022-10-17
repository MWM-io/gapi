package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"github.com/mwm-io/gapi/server/openapi"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/log/cloud_logging"
	"github.com/mwm-io/gapi/server"

	"github.com/mwm-io/gapi/examples/hello-world/internal"
)

func main() {
	ctx := context.Background()
	clientLogger := setupLog(ctx)
	defer clientLogger.Close()

	r := server.NewMux()

	server.AddHandler(r, "GET", "/json/hello", internal.JsonHelloWorldHandlerF())
	server.AddHandler(r, "GET", "/xml/hello", internal.XmlHelloWorldHandlerF())
	server.AddHandler(r, "GET", "/error/hello", internal.ErrorHelloWorldHandlerF())

	server.AddHandler(r, "POST", "/process/hello/{id}", internal.ProcessHandlerF())

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

func setupLog(ctx context.Context) *logging.Client {
	client, err := logging.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Fatalln(err)
	}

	log := gLog.NewDefaultLogger(
		gLog.NewMultiWriter(
			gLog.NewFilterWriter(gLog.InfoSeverity, cloud_logging.NewWriter(client.Logger("application"))),
			gLog.NewWriter(gLog.JSONEntryMarshaler, os.Stdout),
		),
	)

	log = log.WithLabels(map[string]string{"PROJECT_ID": os.Getenv("GOOGLE_CLOUD_PROJECT")})
	gLog.SetGlobalLogger(log)

	return client
}
