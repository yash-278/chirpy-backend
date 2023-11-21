package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yash-278/chirpy-backend/auth"
)

var ErrInvalidCreds = errors.New("invalid credentials")

func (cfg *apiConfig) handlerUserUpdate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	userId, _, err := cfg.UtilAuthRequest(w, r, auth.Issuer.AccessIssuer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	user, err := cfg.DB.GetUser(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var hashedPassword string

	if params.Password != "" {
		hashedPassword, err = auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
			return
		}
	}

	user, err = cfg.DB.UpdateUser(userId, params.Email, hashedPassword, user.IsChirpyRed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.Id,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})
}
