package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var store map[string]entry
var pathIds map[string][]string

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

func notImplemented(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Not implemented", path)

	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Not implemented"))
	check(err)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	e, ok := store[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", e.ContentType)
	_, err := w.Write(e.Value)
	check(err)
}

func errNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, err := w.Write([]byte("Not allowed"))
	check(err)
}

func headHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	value, err := ioutil.ReadAll(r.Body)
	check(err)

	set(path, value, r.Header.Get("Content-Type"))

	getHandler(w, r)
}

// Generates the id for a new element
func idGenerator(path string) (newid string) {
	if val, ok  := pathIds[path]; ok {
		newid = strconv.Itoa(len(val))
	} else {
		newid = "0"
	}
	pathIds[path] = append(pathIds[path], newid)
	return
}	

func postHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	if path[len(path)-1:] == "/" {
		errNotAllowedHandler(w, r)
	} else {
		// generating an id for the new element
		pathParent :=  path[:strings.LastIndex(path, "/")]+"/"
		newid := idGenerator(pathParent)
		
		// generating headers
		w.Header().Add("Location",path+"/"+newid)
		w.WriteHeader(http.StatusCreated)

		mux.Vars(r)["path"] = path+"/"+newid

		putHandler(w, r)
	}
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
	pathIds = make(map[string][]string)
	if overrideContentType != "" {

	}
}
