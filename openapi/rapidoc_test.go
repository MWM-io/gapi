package openapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRapiDocReceiverResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	response, err := NewRapiDocReceiverHandler().Serve(rec, httptest.NewRequest(http.MethodGet, "/oauth-receiver.html", nil))

	if err != nil {
		t.Fatal(err)
	}

	if rec.Header().Get("Content-Type") != "text/html; charset=utf-8" {
		t.Fatal("OAuth callback must be served as HTML")
	}

	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("OAuth callbacks must not be cached")
	}

	body, errRead := io.ReadAll(response.(io.Reader))

	if errRead != nil {
		t.Fatal(errRead)
	}

	if strings.Contains(string(body), "https://unpkg.com/") || strings.Contains(string(body), "<oauth-receiver>") {
		t.Fatal("OAuth callback must not load the upstream receiver that calls window.opener before BroadcastChannel")
	}
}
