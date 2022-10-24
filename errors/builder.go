package errors

import (
	"sync"
)

func init() {
	AddBuilders(defaultBuilder())
}

var (
	builderMux sync.RWMutex
	builders   []ErrorBuilder
)

// AddBuilders will register the given ErrorBuilder to be used.
func AddBuilders(errorBuilders ...ErrorBuilder) {
	builderMux.Lock()
	defer builderMux.Unlock()

	builders = append(builders, errorBuilders...)
}

// Build will call all registered builders to add additional information on the given Error,
// depending on the sourceError
func Build(err Error, sourceError error) Error {
	builderMux.RLock()
	defer builderMux.RUnlock()

	for _, builder := range builders {
		err = builder.Build(err, sourceError)
	}

	return err
}

// ErrorBuilder will try to interpret the sourceErr to populate Error with additional data.
type ErrorBuilder interface {
	Build(err Error, sourceError error) Error
}

// ErrorBuilderFunc is a function that implements the ErrorBuilder interface.
type ErrorBuilderFunc func(err Error, sourceError error) Error

// Build implements the ErrorBuilder interface.
func (e ErrorBuilderFunc) Build(err Error, sourceError error) Error {
	return e(err, sourceError)
}

func defaultBuilder() ErrorBuilder {
	return ErrorBuilderFunc(func(err Error, sourceError error) Error {
		sourceErr, ok := sourceError.(Error)
		if !ok {
			return err
		}

		return err.
			WithStatus(sourceErr.StatusCode()).
			WithSeverity(sourceErr.Severity()).
			WithTimestamp(sourceErr.Timestamp()).
			WithKind(sourceErr.Kind()).
			WithStackTrace(sourceErr.StackTrace())
	})
}
