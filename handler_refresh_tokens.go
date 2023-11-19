package main

import (
	"net/http"
	"strconv"

	"github.com/yash-278/chirpy-backend/auth"
)

func (cfg *apiConfig) handleRefreshTokenCreate(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
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
	if sub, err := claims.GetIssuer(); sub != auth.Issuer.RefreshIssuer {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check for revocations
	if _, err := cfg.DB.GetRefreshTokenById(token); err == nil {
		respondWithError(w, http.StatusUnauthorized, "Token Revoked")
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

	jwtKey, err := auth.CreateJWT(cfg.secretKey, auth.Issuer.AccessIssuer, user.Id, auth.AccessTokenExpireTime)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: jwtKey,
	})
}

func (cfg *apiConfig) handleRefreshTokenRevoke(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
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
	if sub, err := claims.GetIssuer(); sub != auth.Issuer.RefreshIssuer {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Check for revocations
	cfg.DB.RevokeRefreshToken(token)

	respondWithJSON(w, http.StatusOK, response{
		Status: "OK",
	})
}
