package openapi

// OperationDescriptor is able to describe itself as an openapi3 operation
type OperationDescriptor interface {
	Doc(builder *OperationBuilder) error
}
