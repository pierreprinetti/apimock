package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pierreprinetti/apimock/store"
)

type router interface {
	Get(string) (http.Handler, bool)
	Set(string, *http.Request) error
	Del(string) bool
}

func getHandler(resources router) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.String()
		e, ok := resources.Get(path)

		if !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		e.ServeHTTP(rw, req)
	}
}

func putHandler(resources router) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.String()
		if err := resources.Set(path, req); err != nil {
			log.Panic(err)
		}

		e, _ := resources.Get(path)

		e.ServeHTTP(rw, req)
	}
}

func deleteHandler(resources router) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.String()

		if ok := resources.Del(path); !ok {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusNoContent)
	}
}

func optionsHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusNoContent)
}

func newRouter(resources router) http.Handler {
	get := getHandler(resources)
	put := putHandler(resources)
	del := deleteHandler(resources)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			get(rw, req)
		case http.MethodPut:
			put(rw, req)
		case http.MethodDelete:
			del(rw, req)
		case http.MethodOptions:
			optionsHandler(rw, req)
		default:
			msg := fmt.Sprintf("HTTP %s handler not implemented.", req.Method)
			log.Println(msg)
			http.Error(rw, msg, http.StatusNotImplemented)
		}
	})
}

func main() {
	resources := store.New(
		store.WithDefaultContentType(getenv("DEFAULT_CONTENT_TYPE", "text/plain")),
		store.WithContentTypeOverride(getenv("FORCED_CONTENT_TYPE", "")),
	)

	apimock := newRouter(resources)

	withCorsHeaders := newCors(apimock)
	withLogging := newLogger(withCorsHeaders)

	if err := http.ListenAndServe(getenv("HOST", ":80"), withLogging); err != nil {
		log.Fatal(err)
	}
}
