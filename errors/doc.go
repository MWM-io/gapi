/*
Package errors provides a generic error carrying useful information:

 - StackTrace for debugging
 - 2 different messages: one for the end user and one for the developer
 - a Kind: a string that can act as an ID for errors
 - a httpStatus
 - a timestamp of when the error was created

This error implement the json.Marshaler and xml.Marshaler interface
so you can return this error in your http handlers.

It also implements the errors.Unwrap interface, allowing you to get the previous error.

It depends on the gapi/log package,
as the error will implement several of its interfaces to log the errors correctly.

## Create a new error

To create a new error, simply call the Err() function and use the setter to set any data you want.

The setter functions will modify the data and won't create a new error, you don't need to reassign the error.

```go
import (
	"net/http"
	"github.com/mwm-io/gapi/errors"
)

err := errors.Err("my error).
	WithKind("not_found").
	WithStatusCode(http.StatusNotFound).
	WithMessage("not found")
```

## Wrap an existing error

To wrap an existing error, simply call the Wrap() function. You can then modify the new error as you want.

Wrapping a nil error will return a nil value.

```go
import (
	"fmt"
	"github.com/mwm-io/gapi/errors"
)

err := fmt.Errorf("source error")
newErr := errors.Wrap(err, "error").WithKind("new_kind")
```

### Populate data from the source error

You might want to carry more than just the message type from the source error.

In order to do that you need to implement the ErrorBuilder interface
and register your builder with an init function.

Your builder should concerned a single type of error.

You can read an example with the errors/google package.

```go
import (
	"github.com/mwm-io/gapi/errors"
)

func init() {
	errors.AddBuilder()
}

var GrpcCodeErrorBuilder = errors.ErrorBuilderFunc(func(err errors.Error, sourceError error) errors.Error {
	sourceErrI, ok := sourceErr.(interface{ WithKind() string })
	if !ok {
		return err
	}

	return err.WithKind(sourceErrI.WithKind())
})
```

## Why is it an interface ?

You may wonder why we use an Error interface
and why don't we use the FullError struct directly.

This is because `(*FullError)(nil) != nil`: a nil value with a concrete type won't match nil.

In order to keep the idiomatic `if Err("my err") != nil`, we always return the Error interface.

*/
package errors
