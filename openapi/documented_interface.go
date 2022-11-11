package openapi

// Documented is able to describe itself as an openapi3 operation
type Documented interface {
	Doc(builder *DocBuilder) error
}
