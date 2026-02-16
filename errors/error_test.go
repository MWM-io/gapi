package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	sourceErr := fmt.Errorf("original error")
	err := Wrap(sourceErr)

	assert.Equal(t, sourceErr.Error(), err.Error())
	assert.Equal(t, sourceErr.Error(), err.Message())

	wrappedErr := errors.Unwrap(err)
	assert.Equal(t, wrappedErr, sourceErr)

	var nilErr error
	wrappedNilErr := Wrap(nilErr)
	assert.Equal(t, nil, wrappedNilErr)

}

func TestErr(t *testing.T) {
	expectedKind := "errorKind"
	expectedMessage := "this is my error"
	expectedStatusCode := http.StatusInternalServerError

	err := Err(expectedKind, expectedMessage).
		WithStatus(expectedStatusCode)

	assert.Equal(t, expectedKind, err.Kind())
	assert.Equal(t, expectedMessage, err.Message())
	assert.Equal(t, expectedStatusCode, err.StatusCode())

	otherError := errors.New("this is an another error")

	err = err.WithError(otherError)

	assert.Equal(t, expectedMessage, err.Error(), "WithError should not override an explicit error message")
	assert.Equal(t, otherError, errors.Unwrap(err), "WithError should set the source error for debug/logging")
}
