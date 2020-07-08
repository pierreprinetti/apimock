package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs HTTP call information.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
	next http.Handler
}

// newLogger returns a new Logger instance
func newLogger(next http.Handler) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "[apimock] ", 3),
		next:   next,
	}
}

type responseWriterRecorder struct {
	http.ResponseWriter

	status int
}

func (rr *responseWriterRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseWriterRecorder) Write(b []byte) (int, error) {
	if rr.status == 0 {
		rr.status = 200
	}
	return rr.ResponseWriter.Write(b)
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()

	rr := &responseWriterRecorder{ResponseWriter: rw}

	l.next.ServeHTTP(rr, req)

	address := req.RemoteAddr
	if realIP := req.Header.Get("X-REAL-IP"); realIP != "" {
		address = realIP
	}

	l.Printf("%s %s %d %s in %v (%s)", req.Method, req.URL.Path, rr.status, http.StatusText(rr.status), time.Since(start), address)
}
