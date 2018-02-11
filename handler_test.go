package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testrouter struct {
	path string
	body []byte
}

func (tr *testrouter) Get(_ string) (http.Handler, bool) {
	if len(tr.body) == 0 {
		return nil, false
	}
	h := func(rw http.ResponseWriter, _ *http.Request) {
		rw.Write(tr.body)
	}
	return http.HandlerFunc(h), true
}

func (tr *testrouter) Set(path string, req *http.Request) error {
	var err error
	tr.body, err = ioutil.ReadAll(req.Body)
	tr.path = path
	return err
}

func (tr *testrouter) Del(_ string) bool { return true }

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

	handlerWithBody := func(body string) *testrouter {
		return &testrouter{body: []byte(body)}
	}

	tests := [...]struct {
		name   string
		store  *testrouter
		checks []checkFunc
	}{
		{
			"gets",
			handlerWithBody("hey"),
			check(
				hasStatus(200),
				hasContents("hey"),
			),
		},
		{
			"miss is 404",
			&testrouter{},
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

func TestPutHandler(t *testing.T) {
	type checkFunc func(*testrouter, *httptest.ResponseRecorder) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	responseHasStatus := func(want int) checkFunc {
		return func(_ *testrouter, rec *httptest.ResponseRecorder) error {
			if rec.Code != want {
				return fmt.Errorf("expected status %d, found %d", want, rec.Code)
			}
			return nil
		}
	}
	responseHasContents := func(want string) checkFunc {
		return func(_ *testrouter, rec *httptest.ResponseRecorder) error {
			if have := rec.Body.String(); have != want {
				return fmt.Errorf("expected body %q, found %q", want, have)
			}
			return nil
		}
	}
	storeHasPath := func(want string) checkFunc {
		return func(router *testrouter, _ *httptest.ResponseRecorder) error {
			if have := router.path; have != want {
				return fmt.Errorf("expected new path %q, found %q", want, have)
			}
			return nil
		}
	}
	storeHasBody := func(want string) checkFunc {
		return func(router *testrouter, _ *httptest.ResponseRecorder) error {
			if have := string(router.body); have != want {
				return fmt.Errorf("expected new body %q, found %q", want, have)
			}
			return nil
		}
	}

	tests := [...]struct {
		name   string
		path   string
		body   string
		checks []checkFunc
	}{
		{
			"stores a new entry",
			"/wow",
			`{"content": "NEW!"}`,
			check(
				storeHasPath("/wow"),
				storeHasBody(`{"content": "NEW!"}`),
			),
		},
		{
			"returns the newly created entry",
			"/wow",
			`{"content": "NEW!"}`,
			check(
				responseHasStatus(200),
				responseHasContents(`{"content": "NEW!"}`),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.path, strings.NewReader(tc.body))
			store := &testrouter{}
			h := putHandler(store)
			rec := httptest.NewRecorder()
			h(rec, req)
			for _, check := range tc.checks {
				if err := check(store, rec); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
