package main

import (
	//"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var store map[string]entry

type entry struct {
	Value       []byte
	ContentType string
}

type pathMessage struct {
	EndpointName string
	Items        []string
}

var resourcesMap map[string][]string

// returns the path as used inside the program
func getPath(r *http.Request) string {
	//path := mux.Vars(r)["path"]
	path := r.URL.Path
	return strings.Trim(path, "/")
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
func idGenerator(path string) (newId string) {
	/*entry := *new(pathEntry)
	if val, ok := pathIds[path]; ok {
		entry = val
	}
	entry.LastId = entry.LastId + 1
	newid = strconv.Itoa(entry.LastId)
	entry.Value = append(pathIds[path].Value, newid)
	pathIds[path] = entry
	log.Print(path, pathIds[path])*/
	newId = strconv.Itoa(len(resourcesMap[path]))
	resourcesMap[path] = append(resourcesMap[path], newId)
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

func errHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err2 := w.Write([]byte(err.Error()))
	check(err2)
}

func getSuccessHandler(w http.ResponseWriter, e entry) {
	w.Header().Set("Content-Type", e.ContentType)
	_, err := w.Write(e.Value)
	check(err)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := getPath(r)
	e, ok := store[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	getSuccessHandler(w, e)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	path := getPath(r)
	value := checkBody(r)
	
	// generating an id and url for the new element
	newid := idGenerator(path)
	url := path+"/"+newid

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
	// modifies a specific element
	path := mux.Vars(r)["path"]
	value, err := ioutil.ReadAll(r.Body)
	check(err)

	set(path, value, r.Header.Get("Content-Type"))

	getHandler(w, r)
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
	resourcesMap = make(map[string][]string)
	if overrideContentType != "" {

	}
}
