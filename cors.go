package main

import "net/http"

// Cors is a middleware handler that logs the request as it goes in and the response as it goes out.
type Cors struct{}

// NewCors returns a new Cors instance
func newCors() *Cors {
	return new(Cors)
}

func (l *Cors) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "1728000") // Pre-flight info is valid for 20 days
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,X-Api-Key,If-Modified-Since,Cache-Control,Content-Type")
	next(w, r)
}
