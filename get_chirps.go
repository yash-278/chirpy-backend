package main

import (
	"net/http"
	"sort"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, 200, chirps)
}
