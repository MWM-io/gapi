package request

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

func handleHandlerResponse(wr WrappedRequest, result interface{}) {
	if result == nil {
		wr.Response.WriteHeader(http.StatusNoContent)
		return
	}

	var err error
	var b []byte
	switch v := result.(type) {
	case io.ReadCloser:
		defer v.Close()
		_, err = io.Copy(wr.Response, v)

	case io.Reader:
		_, err = io.Copy(wr.Response, v)

	case []byte:
		_, err = wr.Response.Write(v)

	default:
		mt, _, errContent := mime.ParseMediaType(wr.ContentType.String())
		if errContent != nil {
			fmt.Fprintf(os.Stderr, "Warn: ParseMediaType failed %s", errContent.Error())
			http.Error(wr.Response, "malformed Content-Type header", http.StatusBadRequest)
			return
		}

		switch mt {
		case "application/xml":
			b, err = xml.MarshalIndent(result, "", "	")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warn: MarshalIndent failed %s", err.Error())
				http.Error(wr.Response, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = wr.Response.Write(b)

		case "application/json":
			err = json.NewEncoder(wr.Response).Encode(result)

		default:
			http.Error(wr.Response, "unsupported Content-Type header", http.StatusHTTPVersionNotSupported)
			return
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Warn: Write response failed %s", err.Error())
		http.Error(wr.Response, err.Error(), http.StatusInternalServerError)
	}
}
