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

func AddBuilders(errorBuilders ...ErrorBuilder) {
	builderMux.Lock()
	defer builderMux.Unlock()

	builders = append(builders, errorBuilders...)
}

func Build(err ErrorI, sourceError error) ErrorI {
	builderMux.RLock()
	defer builderMux.RUnlock()

	for _, builder := range builders {
		err = builder.Build(err, sourceError)
	}

	return err
}

// ErrorBuilder will try to interpret the sourceErr to populate Error with additional data.
type ErrorBuilder interface {
	Build(err ErrorI, sourceError error) ErrorI
}

type ErrorBuilderFunc func(err ErrorI, sourceError error) ErrorI

func (e ErrorBuilderFunc) Build(err ErrorI, sourceError error) ErrorI {
	return e(err, sourceError)
}

func defaultBuilder() ErrorBuilder {
	return ErrorBuilderFunc(func(err ErrorI, sourceError error) ErrorI {
		sourceErr, ok := sourceError.(ErrorI)
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
