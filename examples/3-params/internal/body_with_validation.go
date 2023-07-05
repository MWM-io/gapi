package internal

import (
	"net/http"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/middleware"
)

// MakeBodyWithValidationHandler decode the given body, store it in h.body and reply with the decoded boyd
func MakeBodyWithValidationHandler() handler.Handler {
	h := bodyWithValidationHandler{}
	h.MiddlewareList = []handler.Middleware{
		middleware.JsonBody(&h.body),
		// You can use middleware.Body(&h.body) if you want to parse body with other content type
	}

	return &h
}

type bodyWithValidationHandler struct {
	handler.WithMiddlewares

	body OrderBody
}

// OrderBody : struct that represent the body used by bodyWithValidationHandler
type OrderBody struct {
	// The required tag validates that the value is not the data types default zero value
	// For numbers ensures value is not zero. For strings ensures value is not "".
	// For slices, maps, pointers, interfaces, channels and functions ensures the value is not nil.
	ProductID string `json:"product_id" required:"true"`
	// if you want to make a custom validation you can implement Validate() as the example bellow for quantity.
	Quantity int `json:"quantity"`
}

// Validate is an implementation of middleware.BodyValidation. If user UserBody is given to
// middleware.Body, this function was called automatically and error handled
func (b *OrderBody) Validate() error {
	if b.Quantity < 5 || b.Quantity > 10 {
		return errors.BadRequest("invalid_quantity", "quantity must be between 5 & 10")
	}

	return nil
}

// Serve implements handler.Handler and is the function called when a request is handled
func (h bodyWithValidationHandler) Serve(_ http.ResponseWriter, _ *http.Request) (interface{}, error) {
	return h.body, nil
}
