package middleware

import (
	"encoding/json"
	"encoding/xml"
)

// Decoder is able to unmarshal the body into a value.
type Decoder interface {
	Unmarshal(b []byte, v interface{}) error
}

// FuncAsDecoder is function type that implements Decoder interface.
type FuncAsDecoder func(b []byte, v interface{}) error

// Unmarshal implements the Decoder interface.
func (f FuncAsDecoder) Unmarshal(b []byte, v interface{}) error {
	return f(b, v)
}

// DecoderByContentType contain all built in decoders
var DecoderByContentType = map[string]Decoder{
	"application/json": FuncAsDecoder(json.Unmarshal),
	"application/xml":  FuncAsDecoder(xml.Unmarshal),
	// TODO url encoded
	// TODO Form data
}
