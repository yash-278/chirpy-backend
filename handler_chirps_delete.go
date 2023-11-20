package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yash-278/chirpy-backend/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	userId, _, err := cfg.UtilAuthRequest(w, r, auth.Issuer.AccessIssuer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "only author can delete chirps")
		return
	}

	cfg.DB.DeleteChirpById(chirp.Id)

	respondWithJSON(w, http.StatusOK, "Deleted")
}
