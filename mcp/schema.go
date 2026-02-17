package mcp

import (
	"encoding/json"
	"reflect"
	"strings"
)

// ParamSource tracks where an MCP tool argument should be routed when executing the handler.
type ParamSource struct {
	Source string // "path", "query", or "body"
	Field  string // original field name in the source struct
}

// schemaProperty represents a single property in a JSON Schema object.
type schemaProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Pattern     string      `json:"pattern,omitempty"`
	Format      string      `json:"format,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

// inputSchema represents a JSON Schema "object" type used as an MCP tool's inputSchema.
type inputSchema struct {
	Type       string                    `json:"type"`
	Properties map[string]schemaProperty `json:"properties,omitempty"`
	Required   []string                  `json:"required,omitempty"`
}

// GenerateInputSchema builds a flat JSON Schema object and a ParamSource map
// from the given parameter and body struct pointers.
// All path/query params and body fields are merged into a single flat schema,
// making the tool input LLM-friendly. The returned paramSourceMap tracks the
// origin of each property so the executor can decompose arguments back into
// path vars, query vars, and body fields.
func GenerateInputSchema(paramPtrs []interface{}, bodyPtr interface{}) (json.RawMessage, map[string]ParamSource, error) {
	schema := inputSchema{
		Type:       "object",
		Properties: make(map[string]schemaProperty),
	}
	sourceMap := make(map[string]ParamSource)

	// Process parameter structs (path and query parameters)
	for _, paramPtr := range paramPtrs {
		if paramPtr == nil {
			continue
		}

		t := reflect.TypeOf(paramPtr)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			continue
		}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			// Check for path tag
			if pathTag := field.Tag.Get("path"); pathTag != "" {
				prop := structFieldToSchemaProperty(field)
				schema.Properties[pathTag] = prop
				schema.Required = append(schema.Required, pathTag)
				sourceMap[pathTag] = ParamSource{Source: "path", Field: field.Name}
				continue
			}

			// Check for query tag (gorilla/schema uses "query" or "schema")
			queryTag := field.Tag.Get("query")
			if queryTag == "" {
				queryTag = field.Tag.Get("schema")
			}
			if queryTag != "" {
				prop := structFieldToSchemaProperty(field)
				schema.Properties[queryTag] = prop
				if field.Tag.Get("required") == "true" {
					schema.Required = append(schema.Required, queryTag)
				}
				sourceMap[queryTag] = ParamSource{Source: "query", Field: field.Name}
			}
		}
	}

	// Process body struct
	if bodyPtr != nil {
		t := reflect.TypeOf(bodyPtr)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			processBodyStruct(t, "", &schema, sourceMap)
		}
	}

	raw, err := json.Marshal(schema)
	if err != nil {
		return nil, nil, err
	}

	return raw, sourceMap, nil
}

// processBodyStruct extracts JSON Schema properties from a body struct type.
func processBodyStruct(t reflect.Type, prefix string, schema *inputSchema, sourceMap map[string]ParamSource) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		fieldName := field.Name
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		schemaKey := fieldName
		if prefix != "" {
			schemaKey = prefix + "." + fieldName
		}

		// For nested structs, flatten them with dot notation
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		if ft.Kind() == reflect.Struct && ft != reflect.TypeOf(json.RawMessage{}) {
			processBodyStruct(ft, schemaKey, schema, sourceMap)
			continue
		}

		prop := structFieldToSchemaProperty(field)
		schema.Properties[schemaKey] = prop
		sourceMap[schemaKey] = ParamSource{Source: "body", Field: fieldName}

		if field.Tag.Get("required") == "true" {
			schema.Required = append(schema.Required, schemaKey)
		}
	}
}

// structFieldToSchemaProperty converts a reflect.StructField to a JSON Schema property.
func structFieldToSchemaProperty(field reflect.StructField) schemaProperty {
	prop := schemaProperty{
		Type: goTypeToJSONSchemaType(field.Type),
	}

	if pattern := field.Tag.Get("pattern"); pattern != "" {
		prop.Pattern = pattern
	}

	if enumStr := field.Tag.Get("enum"); enumStr != "" {
		prop.Enum = strings.Split(enumStr, ",")
	}

	if desc := field.Tag.Get("description"); desc != "" {
		prop.Description = desc
	}

	return prop
}

// goTypeToJSONSchemaType maps a Go type to a JSON Schema type string.
func goTypeToJSONSchemaType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "string"
	}
}
