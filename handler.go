package main

import (
	"log"
	"net/http"
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
