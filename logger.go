package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/negroni"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// newLogger returns a new Logger instance
func newLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[apimock] ", 3)}
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	// l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(rw, req)

	res := rw.(negroni.ResponseWriter)
	address := req.Header.Get("X-REAL-IP")
	if address == "" {
		address = req.RemoteAddr
	}
	l.Printf("%s %s %v %s in %v (%v)", req.Method, req.URL.Path, res.Status(), http.StatusText(res.Status()), time.Since(start), address)
}
