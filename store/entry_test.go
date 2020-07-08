package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContentTypeFromRequest(t *testing.T) {
	type checkFunc func(string) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasValue := func(want string) checkFunc {
		return func(have string) error {
			if have != want {
				return fmt.Errorf("expected Content-Type %q, found %q", want, have)
			}
			return nil
		}
	}

	testCases := [...]struct {
		name        string
		contentType string
		override    string
		def         string
		checks      []checkFunc
	}{
		{
			name:        "parses the content-type from req",
			contentType: "whatever",
			checks:      check(hasValue("whatever")),
		},
		{
			name:        "uses the content-type from req if present",
			contentType: "whatever",
			def:         "content/default",
			checks:      check(hasValue("whatever")),
		},
		{
			name:   "uses default when absent from req",
			def:    "content/default",
			checks: check(hasValue("content/default")),
		},
		{
			name:     "uses override value if present",
			override: "the/override",
			checks:   check(hasValue("the/override")),
		},
		{
			name:        "uses override value over current",
			contentType: "ignore/this",
			override:    "the/override",
			checks:      check(hasValue("the/override")),
		},
		{
			name:     "uses override value over default",
			override: "the/override",
			def:      "content/default",
			checks:   check(hasValue("the/override")),
		},
		{
			name:        "uses override value over anything",
			contentType: "ignore/this",
			override:    "the/override",
			def:         "content/default",
			checks:      check(hasValue("the/override")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatalf("generating the request: %v", err)
			}
			req.Header.Add("Content-Type", tc.contentType)

			ct := contentTypeFromRequest(req, tc.override, tc.def)

			for _, check := range tc.checks {
				if err := check(ct); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestEntryServeHTTP(t *testing.T) {
	type checkFunc func(*httptest.ResponseRecorder) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasContentType := func(want string) checkFunc {
		return func(rw *httptest.ResponseRecorder) error {
			if have := rw.Result().Header.Get("Content-Type"); have != want {
				return fmt.Errorf("expected content-type %q, found %q", want, have)
			}
			return nil
		}
	}
	hasBody := func(want string) checkFunc {
		return func(rw *httptest.ResponseRecorder) error {
			body, err := ioutil.ReadAll(rw.Result().Body)
			if err != nil {
				return fmt.Errorf("reading the body: %v", err)
			}
			rw.Result().Body = ioutil.NopCloser(bytes.NewBuffer(body))
			if have := string(body); have != want {
				return fmt.Errorf("expected body %q, found %q", want, have)
			}
			return nil
		}
	}

	testCases := [...]struct {
		name        string
		contentType string
		body        string
		checks      []checkFunc
	}{
		{
			name:        "sets the content type",
			contentType: "our/contenttype",
			checks:      check(hasContentType("our/contenttype")),
		},
		{
			name:   "sends the body",
			body:   "this is the body",
			checks: check(hasBody("this is the body")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := entry{
				contentType: tc.contentType,
				body:        []byte(tc.body),
			}

			req := httptest.NewRequest("GET", "http://example.com/foo", nil)
			rw := httptest.NewRecorder()
			e.ServeHTTP(rw, req)

			for _, check := range tc.checks {
				if err := check(rw); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
