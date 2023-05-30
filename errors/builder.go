package errors

import (
	"sync"
)

var (
	// errorBuilders slice contains all custom ErrorBuilder callbacks.
	// You can update this list if you want to change ErrorBuilder callbacks
	// for all you handlers.
	errorBuilders []ErrorBuilder

	// addErrorBuildersMU is a sync.Mutex used by AddErrorBuilders.
	addErrorBuildersMU sync.Mutex
)

// ErrorBuilder is a callback that transform the given error to a gapi Error.
type ErrorBuilder func(err error) Error

// AddErrorBuilders appends custom errors.ErrorBuilder.
// These callbacks are executed when wrapping an error with errors.Wrap().
func AddErrorBuilders(builders ...ErrorBuilder) {
	addErrorBuildersMU.Lock()
	defer addErrorBuildersMU.Unlock()

	errorBuilders = append(errorBuilders, builders...)
}
