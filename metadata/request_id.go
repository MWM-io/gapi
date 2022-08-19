package metadata

import (
	"net/http"

	"github.com/google/uuid"
)

// GetRequestID return an ID  for given http.Request.
// 	Return X-Request-ID header if exists
// 	else a generated UUID is returned
func GetRequestID(r *http.Request) string {
	// TODO : Google Cloud Trace specific https://cloud.google.com/trace
	// if traceContext := r.Header.Get("X-Cloud-Trace-Context"); traceContext != "" {
	// 	if separator := strings.Index(traceContext, "/"); separator > 0 {
	// 		traceContext = traceContext[:separator]
	// 	}
	// 	return fmt.Sprintf("projects/%s/traces/%s", ProjectID, traceContext),
	// 		traceContext
	// }

	if requestId := r.Header.Get("X-Request-ID"); requestId != "" {
		return requestId
	}

	return uuid.NewString()
}
