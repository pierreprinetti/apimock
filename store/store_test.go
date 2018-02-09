package store

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestStoreGet(t *testing.T) {
	type checkFunc func(http.Handler, bool) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasBody := func(want string) checkFunc {
		return func(handler http.Handler, _ bool) error {
			if have := handler.(entry).body; string(have) != want {
				return fmt.Errorf("expected body %q, found %q", want, string(have))
			}
			return nil
		}
	}
	hasOk := func(want bool) checkFunc {
		return func(_ http.Handler, have bool) error {
			if have != want {
				return fmt.Errorf("expected ok %v, found %v", want, have)
			}
			return nil
		}
	}

	storeWith := func(path, body string) *Store {
		return &Store{entries: map[string]entry{path: entry{body: []byte(body)}}}
	}

	testCases := [...]struct {
		name   string
		store  *Store
		path   string
		checks []checkFunc
	}{
		{
			"gets existing entry",
			storeWith("this path", "this entry"),
			"this path",
			check(
				hasBody("this entry"),
				hasOk(true),
			),
		},
		{
			"ok is false if not found",
			storeWith("this path", "this entry"),
			"that path",
			check(
				hasOk(false),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h, ok := tc.store.Get(tc.path)
			for _, check := range tc.checks {
				if err := check(h, ok); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestStoreSet(t *testing.T) {
	type checkFunc func(*Store, error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasEntry := func(path, want string) checkFunc {
		return func(s *Store, _ error) error {
			e, ok := s.entries[path]
			if !ok {
				return fmt.Errorf("expected entry with path %q", path)
			}
			if have := string(e.body); have != want {
				return fmt.Errorf("expected body %q, found %q", want, have)
			}
			return nil
		}
	}
	hasError := func(want error) checkFunc {
		return func(_ *Store, have error) error {
			if have != want {
				return fmt.Errorf("expected error %v, found %v", want, have)
			}
			return nil
		}
	}

	storeWith := func(path, body string) *Store {
		return &Store{entries: map[string]entry{path: entry{body: []byte(body)}}}
	}

	testCases := [...]struct {
		name    string
		store   *Store
		newPath string
		newBody string
		checks  []checkFunc
	}{
		{
			"adds a new entry",
			storeWith("one path", "one entry"),
			"new path",
			"new entry",
			check(
				hasEntry("new path", "new entry"),
				hasError(nil),
			),
		},
		{
			"overwrites an existing entry",
			storeWith("one path", "one entry"),
			"one path",
			"two entries",
			check(
				hasEntry("one path", "two entries"),
				hasError(nil),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", tc.newPath, strings.NewReader(tc.newBody))
			if err != nil {
				t.Fatalf("creating the request: %v", err)
			}

			e := tc.store.Set(tc.newPath, req)
			for _, check := range tc.checks {
				if err := check(tc.store, e); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestStoreDel(t *testing.T) {
	type checkFunc func(*Store, bool) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasNoEntry := func(path string) checkFunc {
		return func(s *Store, _ bool) error {
			_, ok := s.entries[path]
			if ok {
				return fmt.Errorf("unexpected entry with path %q", path)
			}
			return nil
		}
	}
	hasOk := func(want bool) checkFunc {
		return func(_ *Store, have bool) error {
			if have != want {
				return fmt.Errorf("expected ok %v, found %v", want, have)
			}
			return nil
		}
	}

	storeWith := func(path, body string) *Store {
		return &Store{entries: map[string]entry{path: entry{body: []byte(body)}}}
	}

	testCases := [...]struct {
		name   string
		store  *Store
		path   string
		checks []checkFunc
	}{
		{
			"removes an existing entry",
			storeWith("one path", "one entry"),
			"one path",
			check(
				hasNoEntry("one path"),
				hasOk(true),
			),
		},
		{
			"signals entry not found",
			storeWith("one path", "one entry"),
			"another path",
			check(
				hasOk(false),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ok := tc.store.Del(tc.path)
			for _, check := range tc.checks {
				if err := check(tc.store, ok); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
