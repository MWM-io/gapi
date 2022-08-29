package request

import (
	"bytes"
	"io"
	"net/http"
)

// ResponseWriter This ResponseWriter is used to store the statusCode
// and the content written to the http.ResponseWriter.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	content    *bytes.Buffer
}

// NewResponseWriter Create a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{
		ResponseWriter: w,
		statusCode:     0,
		content:        new(bytes.Buffer),
	}
}

// StatusCode return the statusCode of the response.
// It's 0 if it isn't set.
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Content returns the content already written to the response.
func (rw *ResponseWriter) Content() io.Reader {
	return rw.content
}

// WriteHeader Write the code in local and to the http response
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the connection as part of an HTTP reply.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	rw.content.Write(b)
	return rw.ResponseWriter.Write(b)
}
