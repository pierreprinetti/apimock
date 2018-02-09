package store

import (
	"io/ioutil"
	"net/http"
	"sync"
)

type Store struct {
	sync.RWMutex
	entries map[string]entry

	overrideContentType string
	defaultContentType  string
}

func (s *Store) Get(path string) (http.Handler, bool) {
	s.RLock()
	defer s.RUnlock()

	e, ok := s.entries[path]
	return http.Handler(e), ok
}

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

func WithDefaultContentType(defaultContentType string) option {
	return func(s *Store) {
		s.defaultContentType = defaultContentType
	}
}

func WithContentTypeOverride(overrideContentType string) option {
	return func(s *Store) {
		s.overrideContentType = overrideContentType
	}
}

func New(options ...option) *Store {
	var s Store

	for _, apply := range options {
		apply(&s)
	}

	return &s
}
