package main

import (
	"encoding/json"
	"net/http"
)

type errors struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respBody := errors{
		Error: msg,
	}
	dat, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
