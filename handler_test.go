package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testrouter struct {
	content string
}

func (tr testrouter) Get(_ string) (http.Handler, bool) {
	if tr.content == "" {
		return nil, false
	}

	h := func(rw http.ResponseWriter, _ *http.Request) {
		rw.Write([]byte(tr.content))
	}
	return http.HandlerFunc(h), true
}
func (tr testrouter) Set(_ string, _ *http.Request) error { return nil }
func (tr testrouter) Del(_ string) bool                   { return true }

func TestGetHandler(t *testing.T) {
	type checkFunc func(*httptest.ResponseRecorder) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasStatus := func(want int) checkFunc {
		return func(rec *httptest.ResponseRecorder) error {
			if rec.Code != want {
				return fmt.Errorf("expected status %d, found %d", want, rec.Code)
			}
			return nil
		}
	}
	hasContents := func(want string) checkFunc {
		return func(rec *httptest.ResponseRecorder) error {
			if have := rec.Body.String(); have != want {
				return fmt.Errorf("expected body %q, found %q", want, have)
			}
			return nil
		}
	}

	tests := [...]struct {
		name   string
		store  testrouter
		checks []checkFunc
	}{
		{
			"gets",
			testrouter{content: "hey"},
			check(
				hasStatus(200),
				hasContents("hey"),
			),
		},
		{
			"miss is 404",
			testrouter{},
			check(
				hasStatus(404),
			),
		},
	}

	req, _ := http.NewRequest("GET", "http://foo.com/", nil)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := getHandler(tc.store)
			rec := httptest.NewRecorder()
			h(rec, req)
			for _, check := range tc.checks {
				if err := check(rec); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
