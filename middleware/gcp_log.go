package middleware

import (
	"context"
	"net/http"
	"os"
	"sync"

	"cloud.google.com/go/logging"
	"github.com/google/uuid"

	"github.com/mwm-io/gapi/handler"
	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/log/cloud_logging"
)

var (
	gcpWriter   *cloud_logging.Writer
	gcpWriterMu sync.RWMutex
)

// GCPLog use structured logs as well as Google own logging client.
// Logs entries are grouped using a unique trace.
type GCPLog struct {
	Logger *gLog.Logger
}

// Wrap implements the request.Middleware interface
func (m GCPLog) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		wr, err := writer()
		if err != nil {
			return nil, err
		}

		traceID := uuid.NewString()

		*m.Logger = *gLog.
			NewLogger(wr, gLog.TimestampNowOpt(), gLog.TracingOpt(traceID, "", false)).
			WithContext(r.Context())

		return h.Serve(w, r)
	})
}

// writer return the GCP logger writer.
func writer() (*cloud_logging.Writer, error) {
	gcpWriterMu.RLock()
	if gcpWriter != nil {
		defer gcpWriterMu.RUnlock()

		return gcpWriter, nil
	}

	gcpWriterMu.RUnlock()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	lc, err := logging.NewClient(context.Background(), "projects/"+projectID)
	if err != nil {
		return nil, err
	}

	setWriter(cloud_logging.NewWriter(lc.Logger(projectID), projectID))

	return writer()
}

// setWriter sets gcp logger writer.
func setWriter(w *cloud_logging.Writer) {
	gcpWriterMu.Lock()
	defer gcpWriterMu.Unlock()

	gcpWriter = w
}
