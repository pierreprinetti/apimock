package store

import (
	"io/ioutil"
	"net/http"
	"sync"
)

// Store saves an HTTP request data associated to a string key.
// It is safe for concurrent usage.
// Store is not directly usable; please initialise one with New.
type Store struct {
	sync.RWMutex
	entries map[string]entry

	overrideContentType string
	defaultContentType  string
}

// Get returns the HTTP request data.
// The returned handler will send back the original HTTP request content type and body.
// The returned boolean is true if a request was found associated to the given key string.
func (s *Store) Get(path string) (http.Handler, bool) {
	s.RLock()
	defer s.RUnlock()

	e, ok := s.entries[path]
	return http.Handler(e), ok
}

// Set saves a request's data associated to a key string.
// An error is returned if the request body io.Reader is not readable.
func (s *Store) Set(path string, req *http.Request) error {
	s.Lock()
	defer s.Unlock()

	contentType := contentTypeFromRequest(req, s.overrideContentType, s.defaultContentType)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	s.entries[path] = entry{contentType, body}

	return nil
}

// Del deletes the entry associated with the given key.
// The returned boolean is true if an entry was actually associated to the given key.
func (s *Store) Del(path string) bool {
	s.Lock()
	defer s.Unlock()

	_, ok := s.entries[path]
	if !ok {
		return false
	}

	delete(s.entries, path)

	return true
}

type option func(*Store)

// WithDefaultContentType is a functional option to modify the behaviour of New.
// The string argument will be used as a content-type when the HTTP request to be saved
// doesn't have one, or is set to the empty string.
func WithDefaultContentType(defaultContentType string) option {
	return func(s *Store) {
		s.defaultContentType = defaultContentType
	}
}

// WithContentTypeOverride is a functional option to modify the behaviour of New.
// The string argument will be used as a content-type for every saved HTTP request.
func WithContentTypeOverride(overrideContentType string) option {
	return func(s *Store) {
		s.overrideContentType = overrideContentType
	}
}

// New initialises a new Store.
func New(options ...option) *Store {
	s := Store{
		entries: make(map[string]entry),
	}

	for _, apply := range options {
		apply(&s)
	}

	return &s
}
