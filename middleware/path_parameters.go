package middleware

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/server"
	"github.com/mwm-io/gapi/server/openapi"
)

// PathParameters is a pre-processor that will set the request parameters into the Parameters field.
type PathParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m PathParameters) Wrap(h server.Handler) server.Handler {
	return server.HandlerFunc(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Parameters == nil {
			return h.Serve(w, r)
		}

		v := reflect.Indirect(reflect.ValueOf(m.Parameters))
		typeOfParameters := v.Type()

		for i := 0; i < v.NumField(); i++ {
			pathParam := typeOfParameters.Field(i).Tag.Get("path")
			val, ok := mux.Vars(r)[pathParam]
			if !ok {
				return nil, errors.Errorf(http.StatusInternalServerError, "unknown path params")
			}

			field := v.FieldByName(typeOfParameters.Field(i).Name)
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, errors.Errorf(http.StatusBadRequest, "%s must be a number", typeOfParameters.Field(i).Name)
				}
				field.SetInt(x)

			case reflect.Float64, reflect.Float32:
				x, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errors.Errorf(http.StatusBadRequest, "%s must be a float", typeOfParameters.Field(i).Name)
				}
				field.SetFloat(x)

			case reflect.Bool:
				field.SetBool(val == "true")

			case reflect.Slice:
				if reflect.TypeOf(i) == reflect.TypeOf([]byte(nil)) {
					field.SetBytes([]byte(val))
				} else {
					return nil, errors.Errorf(http.StatusBadRequest, "cannot have a slice in parameters")
				}

			case reflect.String:
				field.SetString(val)
			default:
				return nil, errors.Errorf(http.StatusBadRequest, "cannot have a parameter with %q type", field.Kind().String())
			}
		}

		return h.Serve(w, r)
	})
}

// Doc implements the openapi.OperationDescriptor interface
func (m PathParameters) Doc(builder *openapi.OperationBuilder) error {
	return builder.WithParams(m.Parameters).Error()
}
