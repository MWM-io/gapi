// Package google provide and register error builders to interpret google errors into gapi errors.
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

var GrpcCodeErrorBuilder = errors.ErrorBuilderFunc(func(err errors.ErrorI, sourceError error) errors.ErrorI {

	switch status.Code(err) {

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
