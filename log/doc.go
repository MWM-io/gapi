/*
Package log

Gapi logging is based on zap.Logger and use its own global logger system to print logs.
This global logger can be overridden with your own zap.Logger

## Usage

### Public methods

Gapi logging expose a list of public methods to get logger with various logging levels:
- log.Debug(ctx context.Context)
- log.Info(ctx context.Context)
- log.Warn(ctx context.Context)
- log.Error(ctx context.Context)
- log.Critical(ctx context.Context)
- log.Alert(ctx context.Context)
- log.Emergency(ctx context.Context)

	import (
		"github.com/mwm-io/gapi/handler"
		gLog "github.com/mwm-io/gapi/log"
		"github.com/mwm-io/gapi/server"
	)

	func main() {
		r := server.NewMux()

		server.AddHandler(r, "GET", "/", handler.Func(HelloWorldHandler))

		gLog.Info(context.Background()).LogMsg("Starting http server")
		if err := server.ServeAndHandleShutdown(r); err != nil {
			gLog.Error(context.Background()).LogError(err)
		}

		gLog.Info(context.Background()).LogMsg("Server stopped")
	}

### Global instance

Logger instance can be retrieved using log.Logger().

	log := log.Logger()
	log.Info("my log")

As mention earlier, you can override gapi logger with log.SetLogger() by passing
your custom zap.Logger instance.

	    myLogger := zap.NewProduction()
		log.SetLogger(myLogger)
		log.Log("my log")

### Context

Context can also be used to store/get logger:
- log.FromContext(ctx) [retrieve gapi logger from context]
- log.NewContext(ctx, logger) [store given zap.Logger into context]
*/
package log
