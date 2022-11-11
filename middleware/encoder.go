package middleware

import (
	"encoding/json"
	"encoding/xml"
)

// Encoder is able to marshal.
type Encoder interface {
	Marshal(v interface{}) ([]byte, error)
}

// FuncAsEncoder is function type that implements Encoder interface.
type FuncAsEncoder func(v interface{}) ([]byte, error)

// Marshal implements the Encoder interface.
func (f FuncAsEncoder) Marshal(v interface{}) ([]byte, error) {
	return f(v)
}

// EncoderByContentType contain all built in decoders
var EncoderByContentType = map[string]Encoder{
	"application/json": FuncAsEncoder(json.Marshal),
	"application/xml":  FuncAsEncoder(xml.Marshal),
	// TODO url encode
}
