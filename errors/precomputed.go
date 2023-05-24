package errors

import (
	"net/http"
)

// BadRequest return an error with status code http.StatusBadRequest
func BadRequest(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusBadRequest).WithKind(kind)
}

// Unauthorized return an error with status code http.StatusUnauthorized
func Unauthorized(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusUnauthorized).WithKind(kind)
}

// PaymentRequired return an error with status code http.StatusPaymentRequired
func PaymentRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusPaymentRequired).WithKind(kind)
}

// Forbidden return an error with status code http.StatusForbidden
func Forbidden(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusForbidden).WithKind(kind)
}

// NotFound return an error with status code http.StatusNotFound
func NotFound(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusNotFound).WithKind(kind)
}

// MethodNotAllowed return an error with status code http.StatusMethodNotAllowed
func MethodNotAllowed(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusMethodNotAllowed).WithKind(kind)
}

// NotAcceptable return an error with status code http.StatusNotAcceptable
func NotAcceptable(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusNotAcceptable).WithKind(kind)
}

// ProxyAuthRequired return an error with status code http.StatusProxyAuthRequired
func ProxyAuthRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusProxyAuthRequired).WithKind(kind)
}

// RequestTimeout return an error with status code http.StatusRequestTimeout
func RequestTimeout(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusRequestTimeout).WithKind(kind)
}

// Conflict return an error with status code http.StatusConflict
func Conflict(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusConflict).WithKind(kind)
}

// Gone return an error with status code http.StatusGone
func Gone(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusGone).WithKind(kind)
}

// LengthRequired return an error with status code http.StatusLengthRequired
func LengthRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusLengthRequired).WithKind(kind)
}

// PreconditionFailed return an error with status code http.StatusPreconditionFailed
func PreconditionFailed(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusPreconditionFailed).WithKind(kind)
}

// RequestEntityTooLarge return an error with status code http.StatusRequestEntityTooLarge
func RequestEntityTooLarge(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusRequestEntityTooLarge).WithKind(kind)
}

// RequestURITooLong return an error with status code http.StatusRequestURITooLong
func RequestURITooLong(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusRequestURITooLong).WithKind(kind)
}

// UnsupportedMediaType return an error with status code http.StatusUnsupportedMediaType
func UnsupportedMediaType(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusUnsupportedMediaType).WithKind(kind)
}

// RequestedRangeNotSatisfiable return an error with status code http.StatusRequestedRangeNotSatisfiable
func RequestedRangeNotSatisfiable(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusRequestedRangeNotSatisfiable).WithKind(kind)
}

// ExpectationFailed return an error with status code http.StatusExpectationFailed
func ExpectationFailed(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusExpectationFailed).WithKind(kind)
}

// Teapot return an error with status code http.StatusTeapot
func Teapot(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusTeapot).WithKind(kind)
}

// MisdirectedRequest return an error with status code http.StatusMisdirectedRequest
func MisdirectedRequest(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusMisdirectedRequest).WithKind(kind)
}

// UnprocessableEntity return an error with status code http.StatusUnprocessableEntity
func UnprocessableEntity(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusUnprocessableEntity).WithKind(kind)
}

// Locked return an error with status code http.StatusLocked
func Locked(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusLocked).WithKind(kind)
}

// FailedDependency return an error with status code http.StatusFailedDependency
func FailedDependency(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusFailedDependency).WithKind(kind)
}

// TooEarly return an error with status code http.StatusTooEarly
func TooEarly(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusTooEarly).WithKind(kind)
}

// UpgradeRequired return an error with status code http.StatusUpgradeRequired
func UpgradeRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusUpgradeRequired).WithKind(kind)
}

// PreconditionRequired return an error with status code http.StatusPreconditionRequired
func PreconditionRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusPreconditionRequired).WithKind(kind)
}

// TooManyRequests return an error with status code http.StatusTooManyRequests
func TooManyRequests(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusTooManyRequests).WithKind(kind)
}

// RequestHeaderFieldsTooLarge return an error with status code http.StatusRequestHeaderFieldsTooLarge
func RequestHeaderFieldsTooLarge(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusRequestHeaderFieldsTooLarge).WithKind(kind)
}

// UnavailableForLegalReasons return an error with status code http.StatusUnavailableForLegalReasons
func UnavailableForLegalReasons(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusUnavailableForLegalReasons).WithKind(kind)
}

// InternalServerError return an error with status code http.StatusInternalServerError
func InternalServerError(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusInternalServerError).WithKind(kind)
}

// NotImplemented return an error with status code http.StatusNotImplemented
func NotImplemented(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusNotImplemented).WithKind(kind)
}

// BadGateway return an error with status code http.StatusBadGateway
func BadGateway(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusBadGateway).WithKind(kind)
}

// ServiceUnavailable return an error with status code http.StatusServiceUnavailable
func ServiceUnavailable(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusServiceUnavailable).WithKind(kind)
}

// GatewayTimeout return an error with status code http.StatusGatewayTimeout
func GatewayTimeout(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusGatewayTimeout).WithKind(kind)
}

// HTTPVersionNotSupported return an error with status code http.StatusHTTPVersionNotSupported
func HTTPVersionNotSupported(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusHTTPVersionNotSupported).WithKind(kind)
}

// VariantAlsoNegotiates return an error with status code http.StatusVariantAlsoNegotiates
func VariantAlsoNegotiates(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusVariantAlsoNegotiates).WithKind(kind)
}

// InsufficientStorage return an error with status code http.StatusInsufficientStorage
func InsufficientStorage(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusInsufficientStorage).WithKind(kind)
}

// LoopDetected return an error with status code http.StatusLoopDetected
func LoopDetected(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusLoopDetected).WithKind(kind)
}

// NotExtended return an error with status code http.StatusNotExtended
func NotExtended(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusNotExtended).WithKind(kind)
}

// NetworkAuthenticationRequired return an error with status code http.StatusNetworkAuthenticationRequired
func NetworkAuthenticationRequired(kind, msgFormat string, args ...interface{}) Error {
	return Err(msgFormat, args...).WithStatus(http.StatusNetworkAuthenticationRequired).WithKind(kind)
}
