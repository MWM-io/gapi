package errors

import (
	"net/http"
)

// BadRequest return an error with status code http.StatusBadRequest
func BadRequest(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusBadRequest)
}

// Unauthorized return an error with status code http.StatusUnauthorized
func Unauthorized(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusUnauthorized)
}

// PaymentRequired return an error with status code http.StatusPaymentRequired
func PaymentRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusPaymentRequired)
}

// Forbidden return an error with status code http.StatusForbidden
func Forbidden(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusForbidden)
}

// NotFound return an error with status code http.StatusNotFound
func NotFound(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusNotFound)
}

// MethodNotAllowed return an error with status code http.StatusMethodNotAllowed
func MethodNotAllowed(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusMethodNotAllowed)
}

// NotAcceptable return an error with status code http.StatusNotAcceptable
func NotAcceptable(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusNotAcceptable)
}

// ProxyAuthRequired return an error with status code http.StatusProxyAuthRequired
func ProxyAuthRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusProxyAuthRequired)
}

// RequestTimeout return an error with status code http.StatusRequestTimeout
func RequestTimeout(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusRequestTimeout)
}

// Conflict return an error with status code http.StatusConflict
func Conflict(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusConflict)
}

// Gone return an error with status code http.StatusGone
func Gone(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusGone)
}

// LengthRequired return an error with status code http.StatusLengthRequired
func LengthRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusLengthRequired)
}

// PreconditionFailed return an error with status code http.StatusPreconditionFailed
func PreconditionFailed(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusPreconditionFailed)
}

// RequestEntityTooLarge return an error with status code http.StatusRequestEntityTooLarge
func RequestEntityTooLarge(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusRequestEntityTooLarge)
}

// RequestURITooLong return an error with status code http.StatusRequestURITooLong
func RequestURITooLong(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusRequestURITooLong)
}

// UnsupportedMediaType return an error with status code http.StatusUnsupportedMediaType
func UnsupportedMediaType(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusUnsupportedMediaType)
}

// RequestedRangeNotSatisfiable return an error with status code http.StatusRequestedRangeNotSatisfiable
func RequestedRangeNotSatisfiable(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusRequestedRangeNotSatisfiable)
}

// ExpectationFailed return an error with status code http.StatusExpectationFailed
func ExpectationFailed(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusExpectationFailed)
}

// Teapot return an error with status code http.StatusTeapot
func Teapot(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusTeapot)
}

// MisdirectedRequest return an error with status code http.StatusMisdirectedRequest
func MisdirectedRequest(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusMisdirectedRequest)
}

// UnprocessableEntity return an error with status code http.StatusUnprocessableEntity
func UnprocessableEntity(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusUnprocessableEntity)
}

// Locked return an error with status code http.StatusLocked
func Locked(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusLocked)
}

// FailedDependency return an error with status code http.StatusFailedDependency
func FailedDependency(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusFailedDependency)
}

// TooEarly return an error with status code http.StatusTooEarly
func TooEarly(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusTooEarly)
}

// UpgradeRequired return an error with status code http.StatusUpgradeRequired
func UpgradeRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusUpgradeRequired)
}

// PreconditionRequired return an error with status code http.StatusPreconditionRequired
func PreconditionRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusPreconditionRequired)
}

// TooManyRequests return an error with status code http.StatusTooManyRequests
func TooManyRequests(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusTooManyRequests)
}

// RequestHeaderFieldsTooLarge return an error with status code http.StatusRequestHeaderFieldsTooLarge
func RequestHeaderFieldsTooLarge(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusRequestHeaderFieldsTooLarge)
}

// UnavailableForLegalReasons return an error with status code http.StatusUnavailableForLegalReasons
func UnavailableForLegalReasons(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusUnavailableForLegalReasons)
}

// InternalServerError return an error with status code http.StatusInternalServerError
func InternalServerError(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusInternalServerError)
}

// NotImplemented return an error with status code http.StatusNotImplemented
func NotImplemented(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusNotImplemented)
}

// BadGateway return an error with status code http.StatusBadGateway
func BadGateway(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusBadGateway)
}

// ServiceUnavailable return an error with status code http.StatusServiceUnavailable
func ServiceUnavailable(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusServiceUnavailable)
}

// GatewayTimeout return an error with status code http.StatusGatewayTimeout
func GatewayTimeout(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusGatewayTimeout)
}

// HTTPVersionNotSupported return an error with status code http.StatusHTTPVersionNotSupported
func HTTPVersionNotSupported(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusHTTPVersionNotSupported)
}

// VariantAlsoNegotiates return an error with status code http.StatusVariantAlsoNegotiates
func VariantAlsoNegotiates(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusVariantAlsoNegotiates)
}

// InsufficientStorage return an error with status code http.StatusInsufficientStorage
func InsufficientStorage(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusInsufficientStorage)
}

// LoopDetected return an error with status code http.StatusLoopDetected
func LoopDetected(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusLoopDetected)
}

// NotExtended return an error with status code http.StatusNotExtended
func NotExtended(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusNotExtended)
}

// NetworkAuthenticationRequired return an error with status code http.StatusNetworkAuthenticationRequired
func NetworkAuthenticationRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(kind, msgFormat, args...).WithStatus(http.StatusNetworkAuthenticationRequired)
}
