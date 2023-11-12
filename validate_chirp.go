package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Body string `json:"body"`
}

type success struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	chirp := profaneChirp(params.Body)

	respBody := success{
		CleanedBody: chirp,
	}

	respondWithJSON(w, 200, respBody)
}

func profaneChirp(chirp string) string {

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splitChirp := strings.Split(chirp, " ")

	for i, word := range splitChirp {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			splitChirp[i] = "****"
		}
	}

	safeChirp := strings.Join(splitChirp, " ")

	return safeChirp
}
