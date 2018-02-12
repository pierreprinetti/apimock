package store

import (
	"net/http"
)

type entry struct {
	contentType string
	body        []byte
}

func contentTypeFromRequest(req *http.Request, override, def string) string {
	if override != "" {
		return override
	}

	contentType := req.Header.Get("Content-Type")

	if contentType == "" {
		return def
	}

	return contentType
}

func (e entry) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", e.contentType)
	rw.Write(e.body)
}
