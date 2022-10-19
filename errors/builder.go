package errors

import (
	"sync"
)

var (
	builderMux sync.RWMutex
	builders   = []ErrorBuilder{}
)

func AddBuilders(builders ...ErrorBuilder) {
	builderMux.Lock()
	defer builderMux.Unlock()

	builders = append(builders, builders...)
}

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

type ErrorBuilderFunc func(err Error, sourceError error) Error

func (e ErrorBuilderFunc) Build(err Error, sourceError error) Error {
	return e(err, sourceError)
}
