package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	e, ok := store[path]

	if !ok {
		http.Error(w, "Resource not found.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", e.ContentType)
	_, err := w.Write(e.Value)
	check(err)
}

func headHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented.", http.StatusNotImplemented)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	value, err := ioutil.ReadAll(r.Body)
	check(err)

	set(path, value, r.Header.Get("Content-Type"))

	getHandler(w, r)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not implemented.", http.StatusNotImplemented)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]

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
	router := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	router.HandleFunc("/{path:.*}", getHandler).Methods("GET")
	router.HandleFunc("/{path:.*}", headHandler).Methods("HEAD")
	router.HandleFunc("/{path:.*}", putHandler).Methods("PUT")
	router.HandleFunc("/{path:.*}", postHandler).Methods("POST")
	router.HandleFunc("/{path:.*}", deleteHandler).Methods("DELETE")
	router.HandleFunc("/{path:.*}", optionsHandler).Methods("OPTIONS")

	n := negroni.New(negroni.NewRecovery(), newLogger(), newCors())
	n.UseHandler(router)
	n.Run(host)
}

func init() {
	store = make(map[string]entry)
	if overrideContentType != "" {

	}
}
