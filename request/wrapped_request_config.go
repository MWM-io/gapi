package request

// Config describe all configurable parameters for the wrapped request.
type Config struct {
	ContentType ContentType
}

// DefaultConfig is the default configuration used to init the wrapped request
var DefaultConfig = Config{
	ContentType: ApplicationJSON,
}
