package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type userParams struct {
	// these tags indicate how the keys in the JSON should be mapped to the struct fields
	// the struct fields must be exported (start with a capital letter) if you want them parsed
	Body string `json:"email"`
}

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := userParams{}

	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	newUser, err := cfg.DB.CreateUser(params.Body)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	respondWithJSON(w, 201, newUser)
}
