package error

// Error /
type Error interface {
	error
	Message() string
	StatusCode() int
	Unwrap() error
}
