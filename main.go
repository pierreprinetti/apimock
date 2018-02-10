package main

import (
	//"encoding/json"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var store map[string]entry

type entry struct {
	Value       []byte
	ContentType string
}

// this is the resource that will represent a collection of resources
// and will be just like {0, 1, ... }
type pathMessage struct {
	Resources []int
}

// returns the path as used inside the program
func getPath(r *http.Request, trim bool) (string, bool) {
	path := r.URL.Path
	strippedPath := strings.TrimRight(path, "/")
	if trim {
		return strippedPath, path != r.URL.Path // not so elegant
	}
	return path, path[len(path)-1:] == "/"
}

// set syncronizes the data in the internal store with the given
// object.
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

// checkBody checks the validity of the request's body and in case returns it.
func checkBody(r *http.Request) []byte {
	value, err := ioutil.ReadAll(r.Body)
	check(err)
	return value
}

// idGenerator generates and stores the id for a new element
func idGenerator(path string) (newID int) {
	message, _ := store[path]
	parsedMessage := unmarshalElement(message)
	newID = len(parsedMessage.Resources)
	parsedMessage.Resources = append(parsedMessage.Resources, newID)
	j, _ := json.Marshal(parsedMessage)
	set(path, j, "application/json")
	return
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	log.Println("Not implemented", path)

	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Not implemented"))
	check(err)
}

func notFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func notAllowedHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func errHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err2 := w.Write([]byte(err.Error()))
	check(err2)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path, _ := getPath(r, false)
	e, ok := store[path]
	if !ok {
		notFoundHandler(w)
		return
	}
	w.Header().Set("Content-Type", e.ContentType)
	w.Header().Set("Location", path)
	_, err := w.Write(e.Value)
	check(err)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	path, _ := getPath(r, true)
	value := checkBody(r)
	// generating an id and url for the new element
	newid := idGenerator(path)
	url := path + "/" + strconv.Itoa(newid)

	// generating headers
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusCreated)
	set(url, value, r.Header.Get("Content-Type"))
	w.Write(value)
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
	// Here we are favouring idempotence: we don't create new resources here.
	path, _ := getPath(r, true)
	_, ok := store[path]
	if !ok {
		notFoundHandler(w)
		return
	}
	value, err := ioutil.ReadAll(r.Body)
	check(err)

	set(path, value, r.Header.Get("Content-Type"))

	getHandler(w, r)
}

func unmarshalElement(message entry) (parsedMessage pathMessage) {
	err := json.Unmarshal(message.Value, &parsedMessage)
	if err != nil {
		parsedMessage = pathMessage{nil}
	}
	return parsedMessage
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	path, isCollection := getPath(r, false)
	if isCollection {
		notAllowedHandler(w)
		return
	}
	if _, ok := store[path]; !ok {
		notFoundHandler(w)
		return
	}
	// this should be the case of a single element and not a collection
	s := strings.Split(path, "/")
	basePath := strings.Join(s[:len(s)-1], "/")
	elements := unmarshalElement(store[basePath]).Resources
	id, err := strconv.Atoi(s[len(s)-1])
	if err != nil {
		// this is the case where we can't find the parent of the element
		// that we're deleting. This happens because it is a collection.
		notAllowedHandler(w)
		return
	}
	var newElementsAr [0]int
	newElements := newElementsAr[:]
	for i := 0; i < len(elements); i++ {
		if elements[i] != id {
			newElements = append(newElements, elements[i])
		}
	}
	j, _ := json.Marshal(pathMessage{newElements})
	store[basePath] = entry{j, "application/json"}
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
