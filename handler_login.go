package main

import (
	"encoding/json"
	"net/http"

	"github.com/yash-278/chirpy-backend/auth"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	jwtKey, err := auth.CreateJWT(cfg.secretKey, auth.Issuer.AccessIssuer, user.Id, auth.AccessTokenExpireTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshKey, err := auth.CreateJWT(cfg.secretKey, auth.Issuer.RefreshIssuer, user.Id, auth.RefreshTokenExpireTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:           user.Id,
		Email:        user.Email,
		Token:        jwtKey,
		RefreshToken: refreshKey,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
