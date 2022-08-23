package error

// Error /
type Error interface {
	Message() string
	StatusCode() int
	Unwrap() error
}
