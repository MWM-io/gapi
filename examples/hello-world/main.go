package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/logging"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/log/cloud_logging"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/router"

	"github.com/mwm-io/gapi/examples/hello-world/internal"
)

func main() {
	ctx := context.Background()
	clientLogger := setupLog(ctx)
	defer clientLogger.Close()

	r := router.NewStrictMux()

	request.AddHandler(r, "GET", "/json/hello", internal.JsonHelloWorldHandlerF())
	request.AddHandler(r, "GET", "/xml/hello", internal.XmlHelloWorldHandlerF())
	request.AddHandler(r, "GET", "/error/hello", internal.ErrorHelloWorldHandlerF())

	request.AddHandler(r, "POST", "/process/hello/{id}", internal.ProcessHandlerF())

	gLog.Info("Starting http server")

	if err := router.ServeAndHandleShutdown(ctx, r); err != nil {
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
