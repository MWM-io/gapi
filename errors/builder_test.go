package errors

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/stacktrace"
)

func TestBuild(t *testing.T) {
	now := time.Now()
	st := stacktrace.New()

	sourceErr := Err("error").
		WithStatus(http.StatusExpectationFailed).
		WithSeverity(gLog.CriticalSeverity).
		WithTimestamp(now).
		WithKind("kind").
		WithStackTrace(st)

	err := Build(Err("new error"), sourceErr)

	assert.Equal(t, "new error", err.Message())
	assert.Equal(t, "kind", err.Kind())
	assert.Equal(t, http.StatusExpectationFailed, err.StatusCode())
	assert.Equal(t, gLog.CriticalSeverity, err.Severity())
	assert.Equal(t, now, err.Timestamp())
	assert.Equal(t, st, err.StackTrace())
}
