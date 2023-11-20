package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yash-278/chirpy-backend/auth"
)

func (cfg *apiConfig) addChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	userId, _, err := cfg.UtilAuthRequest(w, r, auth.Issuer.AccessIssuer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	chirpStr, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newChirp, err := cfg.DB.CreateChirp(chirpStr, userId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}
