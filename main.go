package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
)

var store map[string]entry

type entry struct {
	Value       []byte
	ContentType string
}

func set(key string, value []byte, contentType string) {
	if overrideContentType != "" {
		contentType = overrideContentType
	} else {
		if contentType == "" {
			contentType = defaultContentType
		}
	}
	store[key] = entry{Value: value, ContentType: contentType}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	e, ok := store[path]

	if !ok {
		http.Error(w, "Resource not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", e.ContentType)
	_, err := w.Write(e.Value)
	if err != nil {
		log.Panic(err)
	}
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()

	value, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	set(path, value, r.Header.Get("Content-Type"))

	getHandler(w, r)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()

	if _, ok := store[path]; !ok {
		http.Error(w, "Resource not found.", http.StatusNotFound)
		return
	}

	delete(store, path)
	w.WriteHeader(http.StatusNoContent)
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			getHandler(rw, req)
		case http.MethodPut:
			putHandler(rw, req)
		case http.MethodDelete:
			deleteHandler(rw, req)
		case http.MethodOptions:
			optionsHandler(rw, req)
		default:
			msg := fmt.Sprintf("HTTP %s handler not implemented.", req.Method)
			log.Println(msg)
			http.Error(rw, msg, http.StatusNotImplemented)
		}
	})

	n := negroni.New(negroni.NewRecovery(), newLogger(), newCors())
	n.UseHandler(router)
	n.Run(host)
}

func init() {
	store = make(map[string]entry)
	if overrideContentType != "" {

	}
}
