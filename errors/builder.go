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
