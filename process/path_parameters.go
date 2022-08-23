package process

import (
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mwm-io/gapi/request"
	"github.com/mwm-io/gapi/response"
)

// PathParameters is a pre-processor that will set the request parameters into the Parameters field.
type PathParameters struct {
	Parameters interface{}
}

// PreProcess implements the server.PreProcess interface
func (m PathParameters) PreProcess(handler request.Handler, r *request.WrappedRequest) (request.Handler, response.Error) {
	v := reflect.Indirect(reflect.ValueOf(m.Parameters))
	typeOfParameters := v.Type()

	for i := 0; i < v.NumField(); i++ {
		pathParam := typeOfParameters.Field(i).Tag.Get("path")
		val, ok := mux.Vars(r.Request)[pathParam]
		if !ok {
			return nil, response.ErrorInternalServerErrorf("unknown path params")
		}

		field := v.FieldByName(typeOfParameters.Field(i).Name)
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, response.ErrorBadRequestf("%s must be a number", typeOfParameters.Field(i).Name)
			}
			field.SetInt(x)

		case reflect.Float64, reflect.Float32:
			x, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, response.ErrorBadRequestf("%s must be a float", typeOfParameters.Field(i).Name)
			}
			field.SetFloat(x)

		case reflect.Bool:
			field.SetBool(val == "true")

		case reflect.Slice:
			if reflect.TypeOf(i) == reflect.TypeOf([]byte(nil)) {
				field.SetBytes([]byte(val))
			} else {
				return nil, response.ErrorBadRequestf("cannot have a slice in parameters")
			}

		case reflect.String:
			field.SetString(val)
		default:
			return nil, response.ErrorBadRequestf("cannot have a parameter with %q type", field.Kind().String())
		}
	}

	return handler, nil
}
