package main

import "net/http"

// Cors is a middleware handler that adds Cross-Origin-Resource-Sharing headers.
type Cors struct {
	next http.Handler
}

// newCors returns a new Cors instance
func newCors(next http.Handler) Cors {
	return Cors{
		next: next,
	}
}

func (m Cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "1728000") // Pre-flight info is valid for 20 days
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,X-Api-Key,If-Modified-Since,Cache-Control,Content-Type")
	m.next.ServeHTTP(w, r)
}
