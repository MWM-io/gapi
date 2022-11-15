package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"

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

		unmarshaler, err := m.resolveContentType(r)
		if err != nil {
			return nil, errors.Wrap(err).WithMessage("unable to resolve content type").WithStatus(http.StatusBadRequest)
		}

		var buffer bytes.Buffer
		reader := io.TeeReader(r.Body, &buffer)

		defer func() {
			_ = r.Body.Close
		}()

		body, err := io.ReadAll(reader)
		if err != nil {
			return nil, errors.Wrap(err).WithMessage("failed to read body").WithStatus(http.StatusBadRequest)
		}

		if errUnmarshal := unmarshaler.Unmarshal(body, m.BodyPtr); errUnmarshal != nil {
			return nil, errors.Wrap(errUnmarshal).WithMessage("failed to unmarshal body").WithStatus(http.StatusBadRequest)
		}

		r.Body = io.NopCloser(bytes.NewReader(buffer.Bytes()))

		if v, ok := m.BodyPtr.(BodyValidation); !m.SkipValidation && ok {
			if errValidate := v.Validate(); errValidate != nil {
				return nil, errors.Wrap(errValidate).WithMessage("validation failed").WithStatus(http.StatusUnprocessableEntity)
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
			return nil, errors.Err(fmt.Sprintf("no content decoder found for content type %s", m.ForcedContentType))
		}

		return result, nil
	}

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = m.DefaultContentType
	}

	wantedType, _, errContent := mime.ParseMediaType(contentType)
	if errContent != nil {
		return nil, errors.Wrap(errContent).WithMessage(fmt.Sprintf("unknown content-type %s", contentType))
	}

	if wantedType == "" {
		wantedType = m.DefaultContentType
	}

	result, ok := m.Decoders[wantedType]
	if !ok || result == nil {
		return nil, errors.Err(fmt.Sprintf("unsupported content-type %s", wantedType))
	}

	return result, nil
}
