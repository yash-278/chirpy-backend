package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/yash-278/chirpy-backend/database"
)

func addChirp(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(dbKey).(*database.DB)

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

	chirpStr, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	newChirp, err := db.CreateChirp(chirpStr)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	respondWithJSON(w, 201, newChirp)
}

func getChirps(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value(dbKey).(*database.DB)

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, 200, chirps)
}
