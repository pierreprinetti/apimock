package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var store map[string]entry
var pathIds map[string]pathEntry

type entry struct {
	Value       []byte
	ContentType string
}

type pathEntry struct {
	Value		[]string
	LastId		int
}

type pathMessage struct {
	EndpointName	string
	Items			[]string
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

func notFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func errHandler(w http.ResponseWriter, err error){
	w.WriteHeader(http.StatusInternalServerError)
	_, err2 := w.Write([]byte(err.Error()))
	check(err2)
}

func getSuccessHandler(w http.ResponseWriter, e entry){
	w.Header().Set("Content-Type", e.ContentType)
	_, err := w.Write(e.Value)
	check(err)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	if len(path)>0 && path[len(path)-1:] == "/"{
		pthId, ok := pathIds[path]
		if !ok {
			notFoundHandler(w)
		}
		j_m, err := json.Marshal(pathMessage{path, pthId.Value})
		if err != nil {
			errHandler(w,err)
		}
		getSuccessHandler(w, entry{[]byte(j_m), "application/json"})

	} else {
		e, ok := store[path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		getSuccessHandler(w, e)
	}
	notFoundHandler(w)
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

// Generates and stores the id for a new element
func idGenerator(path string) (newid string) {
	entry := *new(pathEntry)
	if val, ok  := pathIds[path]; ok {
		entry = val
	}
	entry.LastId = entry.LastId+1
	newid = strconv.Itoa(entry.LastId)
	entry.Value = append(pathIds[path].Value, newid)
	pathIds[path] = entry
	log.Print(path, pathIds[path])
	return
}	

func postHandler(w http.ResponseWriter, r *http.Request) {
	path := mux.Vars(r)["path"]
	if path[len(path)-1:] == "/" {
		errNotAllowedHandler(w, r)
	} else {
		// generating an id for the new element
		newid := idGenerator(path+"/")
		
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
	pathIds = make(map[string]pathEntry)
	if overrideContentType != "" {

	}
}
