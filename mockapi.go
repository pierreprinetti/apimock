package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var store map[string][]byte

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Not implemented", path)

	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Not implemented"))
	check(err)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	value, ok := store[path]

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, err := w.Write(value)
	check(err)
}

func headHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	value, err := ioutil.ReadAll(r.Body)
	check(err)

	store[path] = value
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(value)
	check(err)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]

	if _, ok := store[path]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(store, path)
	w.WriteHeader(http.StatusNoContent)
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func main() {
	r := mux.NewRouter()
	// r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/{path:.*}", getHandler).Methods("GET")
	r.HandleFunc("/{path:.*}", headHandler).Methods("HEAD")
	r.HandleFunc("/{path:.*}", putHandler).Methods("PUT")
	r.HandleFunc("/{path:.*}", postHandler).Methods("POST")
	r.HandleFunc("/{path:.*}", deleteHandler).Methods("DELETE")
	r.HandleFunc("/{path:.*}", optionsHandler).Methods("OPTIONS")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(host, nil))
}

func init() {
	store = make(map[string][]byte)
}
