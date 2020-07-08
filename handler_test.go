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
	path             string
	body             []byte
	deleteCalledWith string
	deleteBool       bool
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

func (tr *testrouter) Del(path string) bool {
	tr.deleteCalledWith = path
	return tr.deleteBool
}

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

	req, _ := http.NewRequest("GET", "http://foo.com/", strings.NewReader(""))
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
			req, _ := http.NewRequest("PUT", tc.path, strings.NewReader(tc.body))
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

func TestDeleteHandler(t *testing.T) {
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
	deleteCalledWith := func(want string) checkFunc {
		return func(router *testrouter, _ *httptest.ResponseRecorder) error {
			if have := router.deleteCalledWith; have != want {
				return fmt.Errorf("expected Del called with path %q, found %q", want, have)
			}
			return nil
		}
	}

	tests := [...]struct {
		name   string
		path   string
		store  *testrouter
		checks []checkFunc
	}{
		{
			"deletes an entry",
			"/wow",
			&testrouter{deleteBool: true},
			check(
				deleteCalledWith("/wow"),
				responseHasStatus(204),
			),
		},
		{
			"returns 404 for unknown routes",
			"/wow",
			&testrouter{deleteBool: false},
			check(
				deleteCalledWith("/wow"),
				responseHasStatus(404),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", tc.path, strings.NewReader(""))
			h := deleteHandler(tc.store)
			rec := httptest.NewRecorder()
			h(rec, req)
			for _, check := range tc.checks {
				if err := check(tc.store, rec); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestOptionsHandler(t *testing.T) {
	t.Run("returns 204", func(t *testing.T) {
		req, _ := http.NewRequest("OPTIONS", "/", strings.NewReader(""))
		rec := httptest.NewRecorder()
		optionsHandler(rec, req)
		if want, have := 204, rec.Code; want != have {
			t.Errorf("expected status code %d, found %d", want, have)
		}
	})
}
