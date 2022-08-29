/*
Package errors provide an error struct with useful information for a http API.

It contains its http error status, and will implement encoding/json and encoding/Xml Marshaler interface
to render itself.
It also has useful debug features:
 - error wrapping
 - stacktrace (TODO)
 - custom user displayed message (vs internal message for debugging) (TODO)
 - implements log package interface for better logging (TODO)

*/
package errors
