package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testrwcors struct{}

func (trw testrwcors) Header() http.Header         { return http.Header{} }
func (trw testrwcors) Write(b []byte) (int, error) { return len(b), nil }
func (trw testrwcors) WriteHeader(int)             {}

type testhandler int

func (h *testhandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) { *h++ }

func TestNewCors(t *testing.T) {
	t.Run("passed handler is called", func(t *testing.T) {
		var h testhandler
		wrapped := newCors(&h)

		req, _ := http.NewRequest("GET", "/", strings.NewReader(""))
		var rw testrwcors

		for range [134]struct{}{} {
			wrapped.ServeHTTP(rw, req)
		}

		if want, have := 134, int(h); want != have {
			t.Errorf("expected the handler to be called %d times, found %d", want, have)
		}
	})

	t.Run("cors headers are injected", func(t *testing.T) {
		var h testhandler
		wrapped := newCors(&h)

		req, _ := http.NewRequest("GET", "/", strings.NewReader(""))
		rw := httptest.NewRecorder()

		wrapped.ServeHTTP(rw, req)

		for k, v := range map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, PUT, POST, DELETE, HEAD, OPTIONS",
			"Access-Control-Allow-Headers": "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,X-Api-Key,If-Modified-Since,Cache-Control,Content-Type",
		} {
			if want, have := v, rw.HeaderMap.Get(k); want != have {
				t.Errorf("expected header %q to have value %q, found %q", k, want, have)
			}
		}
	})

	t.Run("cors cache ttl is injected in OPTIONS calls", func(t *testing.T) {
		var h testhandler
		wrapped := newCors(&h)

		req, _ := http.NewRequest("OPTIONS", "/", strings.NewReader(""))
		rw := httptest.NewRecorder()

		wrapped.ServeHTTP(rw, req)

		key := "Access-Control-Max-Age"
		if want, have := "1728000", rw.HeaderMap.Get(key); want != have {
			t.Errorf("expected header %q to have value %q, found %q", key, want, have)
		}
	})
}
