// Package google provides and registers error builders to interpret google errors into gapi errors.
//
// In order to use it, you just need to import this package:
// import _ "github.com/mwm-io/gapi/errors/google"
package google

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mwm-io/gapi/errors"
)

func init() {
	errors.AddBuilders(
		GrpcCodeErrorBuilder,
	)
}

// GrpcCodeErrorBuilder is an ErrorBuilderFunc that will check the grpc status of the sourceErr
// and set the corresponding http status.
var GrpcCodeErrorBuilder = errors.ErrorBuilderFunc(func(err errors.Error, sourceError error) errors.Error {
	switch status.Code(sourceError) {
	case codes.InvalidArgument:
		err = err.WithStatus(http.StatusNotAcceptable)

	case codes.DeadlineExceeded:
		err = err.WithStatus(http.StatusRequestTimeout)

	case codes.NotFound:
		err = err.WithStatus(http.StatusNotFound)

	case codes.AlreadyExists:
		err = err.WithStatus(http.StatusConflict)

	case codes.PermissionDenied:
		err = err.WithStatus(http.StatusForbidden)

	default:
	}

	return err
})
