package main

import "net/http"

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("200 OK"))
}
