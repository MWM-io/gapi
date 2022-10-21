package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gLog "github.com/mwm-io/gapi/log"
	"github.com/mwm-io/gapi/stacktrace"
)

func ExampleErr() {
	err := Err("this is my error").
		WithMessage("message for the final user").
		WithKind("error").
		WithStatus(http.StatusBadGateway).
		WithSeverity(gLog.CriticalSeverity).
		WithStackTrace(stacktrace.New()).
		WithTimestamp(time.Now())

	js, _ := err.MarshalJSON()

	fmt.Println(string(js))
	// Output: {"message":"message for the final user","kind":"error"}

}

func ExampleWrap() {
	sourceErr := fmt.Errorf("original error")
	err := Wrap(sourceErr, "wrapped error")

	fmt.Println(err.Error())
	js, _ := err.MarshalJSON()
	fmt.Println(string(js))

	wrappedErr := errors.Unwrap(err)
	fmt.Println(wrappedErr.Error())

	var nilErr error
	wrappedNilErr := Wrap(nilErr, "nil error")
	fmt.Printf("%t\n", wrappedNilErr == nil)

	// Output:
	// wrapped error: original error
	// {"message":"wrapped error","kind":""}
	// original error
	// true
}

func TestError_WithStatus(t *testing.T) {
	err := Err("error")
	assert.Equal(t, gLog.ErrorSeverity, err.Severity())

	err.WithStatus(http.StatusExpectationFailed)
	assert.Equal(t, gLog.WarnSeverity, err.Severity())

	err.WithStatus(http.StatusInternalServerError)
	assert.Equal(t, gLog.ErrorSeverity, err.Severity())

}
