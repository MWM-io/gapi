package response

// Error TODO
type Error interface {
	Message() string
	StatusCode() int
	Origin() error
}
