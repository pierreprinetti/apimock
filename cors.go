package main

import "net/http"

// Cors is a middleware handler that logs the request as it goes in and the response as it goes out.
type Cors struct{}

// NewCors returns a new Cors instance
func newCors() *Cors {
	return new(Cors)
}

func (l *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	next(w, r)
}
