package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pierreprinetti/apimock/store"
)

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

	if err := http.ListenAndServe(
		getenv("HOST", ":"+getenv("PORT", "8800")),
		withLogging); err != nil {
		log.Fatal(err)
	}
}
