package middleware

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
)

// PathParameters is a middleware that will set the request path parameters into the Parameters field.
type PathParameters struct {
	Parameters interface{}
}

// Wrap implements the request.Middleware interface
func (m PathParameters) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.Parameters == nil {
			return h.Serve(w, r)
		}

		v := reflect.Indirect(reflect.ValueOf(m.Parameters))
		typeOfParameters := v.Type()

		for i := 0; i < v.NumField(); i++ {
			pathParam := typeOfParameters.Field(i).Tag.Get("path")
			val, ok := mux.Vars(r)[pathParam]
			if !ok {
				return nil, errors.Err("invalid_param_type", "unknown path params").WithStatus(http.StatusInternalServerError)
			}

			field := v.FieldByName(typeOfParameters.Field(i).Name)
			switch field.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				x, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return nil, errors.BadRequest("invalid_param_type", "%s must be a number", typeOfParameters.Field(i).Name).
						WithError(err)
				}
				field.SetInt(x)

			case reflect.Float64, reflect.Float32:
				x, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errors.BadRequest("invalid_param_type", "%s must be a float", typeOfParameters.Field(i).Name).
						WithError(err)
				}
				field.SetFloat(x)

			case reflect.Bool:
				field.SetBool(val == "true")

			case reflect.Slice:
				if reflect.TypeOf(i) == reflect.TypeOf([]byte(nil)) {
					field.SetBytes([]byte(val))
				} else {
					return nil, errors.PreconditionFailed("invalid_param_type", "path params can't be a slice. Caused by parameter '%s'", pathParam)
				}

			case reflect.String:
				field.SetString(val)
			default:
				return nil, errors.PreconditionFailed("invalid_param_type", "type %q isn't supported. Caused by parameter '%s'", pathParam, field.Kind().String()).WithStatus(http.StatusBadRequest)
			}
		}

		return h.Serve(w, r)
	})
}

// Doc implements the openapi.Documented interface
func (m PathParameters) Doc(builder *openapi.DocBuilder) error {
	if m.Parameters == nil {
		return nil
	}

	return builder.WithParams(m.Parameters).Error()
}
