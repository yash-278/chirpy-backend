package main

import (
	"errors"
	"strings"
)

type parameters struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Body string `json:"body"`
}

func validateChirp(input string) (string, error) {

	if len(input) > 140 {
		return "", errors.New("chirp is too long")
	}

	chirp := profaneChirp(input)

	return chirp, nil
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
