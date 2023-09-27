package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/mwm-io/gapi/errors"
	"github.com/mwm-io/gapi/handler"
	"github.com/mwm-io/gapi/openapi"
)

// Body return a preconfigured BodyDecoder with :
//   - BodyDecoder.DefaultContentType = `application/json`
//   - BodyDecoder.Decoders with all referenced Decoder
func Body(bodyPtr interface{}) BodyDecoder {
	return BodyDecoder{
		BodyPtr:            bodyPtr,
		Decoders:           DecoderByContentType,
		DefaultContentType: "application/json",
		SkipValidation:     false,
	}
}

// JsonBody return a preconfigured BodyDecoder with :
//   - BodyDecoder.ForcedContentType = `application/json`
func JsonBody(bodyPtr interface{}) BodyDecoder {
	return BodyDecoder{
		BodyPtr:           bodyPtr,
		Decoders:          DecoderByContentType,
		ForcedContentType: "application/json",
		SkipValidation:    false,
	}
}

// BodyDecoder is a middleware that will Unmarshal the incoming request body into the Body field.
type BodyDecoder struct {
	// BodyPtr is a pointer to the variable you want to unmarshal your request body into.
	BodyPtr interface{}
	// Decoders is the list the available Decoders by content-type.
	Decoders map[string]Decoder
	// DefaultContentType is the default content-type if the request don't have any.
	DefaultContentType string
	// SkipValidation indicates whether you want to skip body validation or not.
	SkipValidation bool
	// ForcedContentType will always decode a body with this content-type.
	ForcedContentType string
}

// SetDefaultContentType set DefaultContentType and return current instance
func (m BodyDecoder) SetDefaultContentType(contentType string) BodyDecoder {
	m.DefaultContentType = contentType
	return m
}

// SetForcedContentType set ForcedContentType and return current instance
func (m BodyDecoder) SetForcedContentType(contentType string) BodyDecoder {
	m.ForcedContentType = contentType
	return m
}

// SetDecoders set Decoders field and return current instance
func (m BodyDecoder) SetDecoders(decoders map[string]Decoder) BodyDecoder {
	m.Decoders = decoders
	return m
}

// AddDecoder add a Decoder in field Decoders and return current instance
func (m BodyDecoder) AddDecoder(contentType string, decoder Decoder) BodyDecoder {
	m.Decoders[contentType] = decoder
	return m
}

// BodyValidation interface can be implemented to trigger an auto validation by BodyDecoder middleware.
type BodyValidation interface {
	Validate() error
}

// Wrap implements the request.Middleware interface
func (m BodyDecoder) Wrap(h handler.Handler) handler.Handler {
	return handler.Func(func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		if m.BodyPtr == nil {
			return h.Serve(w, r)
		}

		if reflect.ValueOf(m.BodyPtr).Kind() != reflect.Ptr {
			return nil, errors.Err("invalid_body_receiver", "BodyPtr must be a pointer")
		}

		unmarshaler, err := m.resolveContentType(r)
		if err != nil {
			return nil, errors.BadRequest("invalid_content_type", "unable to resolve content type").
				WithError(err)
		}

		var buffer bytes.Buffer
		reader := io.TeeReader(r.Body, &buffer)

		defer func() {
			_ = r.Body.Close
		}()

		body, err := io.ReadAll(reader)
		if err != nil {
			return nil, errors.BadRequest("body_error", "failed to read body").
				WithError(err)
		}

		if errUnmarshal := unmarshaler.Unmarshal(body, m.BodyPtr); errUnmarshal != nil {
			return nil, errors.BadRequest("invalid_body_format", "failed to decode body").
				WithError(errUnmarshal)
		}

		r.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		valOf := reflect.ValueOf(m.BodyPtr).Elem()
		if kind := valOf.Kind(); kind == reflect.Struct {
			val := reflect.Indirect(valOf)

			for i := 0; i < val.NumField(); i++ {
				typeOfParameters := val.Type()
				typeOfFieldI := typeOfParameters.Field(i)

				if pattern := typeOfFieldI.Tag.Get("pattern"); pattern != "" {
					rex, errC := regexp.Compile(pattern)
					if errC != nil {
						return nil, errors.InternalServerError("pattern_must_be_regex", "pattern must contain a regular expression")
					}

					if !val.Field(i).CanInterface() {
						return nil, errors.InternalServerError("interface_failed", "interface cannot be used without panicking")
					}

					fieldValue := val.Field(i).Interface()
					if !rex.MatchString(fmt.Sprintf("%v", fieldValue)) {
						return nil, errors.BadRequest("body_validation_failed", "field %s does not match the required pattern", typeOfFieldI.Name)
					}
				}

				if typeOfFieldI.Tag.Get("required") != "true" {
					continue
				}

				if !val.Field(i).IsZero() {
					continue
				}

				fieldName := typeOfFieldI.Name
				switch jsonTag := typeOfFieldI.Tag.Get("json"); jsonTag {
				case "-":
					return nil, errors.InternalServerError("invalid_config", "field '%s' is required but json tag value is '-'", fieldName)

				case "":
					return nil, errors.BadRequest("missing_param", "field %s is required", fieldName)

				default:
					parts := strings.Split(jsonTag, ",")
					name := parts[0]
					if name == "" {
						name = fieldName
					}

					return nil, errors.BadRequest("missing_param", "field %s is required", name)
				}
			}
		}

		if v, ok := m.BodyPtr.(BodyValidation); !m.SkipValidation && ok {
			if errValidate := v.Validate(); errValidate != nil {
				if castedErr, casted := errValidate.(errors.Error); casted {
					return nil, castedErr
				}

				return nil, errors.BadRequest("invalid_body", errValidate.Error())
			}
		}

		return h.Serve(w, r)
	})
}

// Doc implements the openapi.Documented interface
func (m BodyDecoder) Doc(builder *openapi.DocBuilder) error {
	if m.BodyPtr == nil {
		return nil
	}

	for contentType := range m.Decoders {
		builder.WithBody(m.BodyPtr, openapi.WithMimeType(contentType))
	}

	return builder.Error()
}

func (m BodyDecoder) resolveContentType(r *http.Request) (Decoder, error) {
	if m.ForcedContentType != "" {
		result, ok := m.Decoders[m.ForcedContentType]
		if !ok {
			return nil, errors.UnsupportedMediaType("unsupported_content_type", "no decoder found for content type %s", m.ForcedContentType)
		}

		return result, nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = m.DefaultContentType
	}

	wantedType, _, errContent := mime.ParseMediaType(contentType)
	if errContent != nil {
		return nil, errors.UnsupportedMediaType("unsupported_content_type", "unknown content-type %s", contentType).
			WithError(errContent)
	}

	if wantedType == "" {
		wantedType = m.DefaultContentType
	}

	result, ok := m.Decoders[wantedType]
	if !ok || result == nil {
		return nil, errors.UnsupportedMediaType("unsupported_content_type", "unsupported content-type %s", wantedType)
	}

	return result, nil
}
