package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

	token := GetBearerToken(r)
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, ErrInvalidCreds.Error())
		return
	}

	claims, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if sub, _ := claims.GetIssuer(); sub != auth.Issuer.AccessIssuer {
		respondWithError(w, http.StatusUnauthorized, "Bad Token")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	strId, _ := claims.GetSubject()

	id, err := strconv.Atoi(strId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.GetUser(id)
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

	user, err = cfg.DB.UpdateUser(id, params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.Id,
			Email: user.Email,
		},
	})
}
