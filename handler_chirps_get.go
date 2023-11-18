package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	fmt.Println(idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	chirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	respondWithJSON(w, 200, chirps)
}
